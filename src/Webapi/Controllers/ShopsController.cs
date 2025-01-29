using Azure.Storage.Blobs;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.Data;
using ZiggyCreatures.Caching.Fusion;

namespace Webapi.Controllers
{
    [Route("api/shops")]
    [ApiController]
    public class ShopsController : ControllerBase
    {
        private readonly ApplicationDbContext db;
        private IHostEnvironment he;
        private BlobServiceClient blobServiceClient;
        private IServiceProvider sp;
        private IFusionCache cache;
        public ShopsController(ApplicationDbContext context, IHostEnvironment _he, BlobServiceClient _blobServiceClient, IServiceProvider _sp, IFusionCache _cache)
        {
            db = context;
            he = _he;
            blobServiceClient = _blobServiceClient;
            sp = _sp;
            cache = _cache;
        }

        private IQueryable<long> shopManagers()
        {
            return db.Users.Where(_ => _.ShopManage).Select(_ => _.Id);
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<Shop>>> GetShops()
        {
            var r = await db.Shops
                .Where(_ => _.IsDeleted == false && shopManagers().Contains(_.CreatedByID))
                .AsNoTracking()
                .ToListAsync();
            foreach (var shop in r) UserCache.LoadCreatedBy(shop, db);
            return r;
        }

        [HttpGet("{id}")]
        public async Task<ActionResult<Shop>> GetShop(long id)
        {
            var shop = await db.Shops
                .Where(_ => _.IsDeleted == false && _.ID == id && shopManagers().Contains(_.CreatedByID))
                .AsNoTracking()
                .FirstOrDefaultAsync();

            if (shop == null) return NotFound();

            UserCache.LoadCreatedBy(shop, db);
            return shop;
        }

        [HttpPut("{id}")]
        public async Task<IActionResult> PutShop(long id, Shop shop)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var dbobj = db.Shops.First(_ => _.ID == id && _.IsDeleted == false && shopManagers().Contains(_.CreatedByID));
            if (dbobj == null) return NotFound();

            if (dbobj.CreatedByID != cu.Id) return Unauthorized();

            dbobj.Caption = shop.Caption;
            dbobj.Description = shop.Description;
            dbobj.DeliveryConditions = shop.DeliveryConditions;

            await db.SaveChangesAsync();

            return NoContent();
        }

        [HttpPost]
        public async Task<ActionResult<Shop>> PostShop(Shop shop)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            if (!shopManagers().Contains(cu.Id)) return Unauthorized();

            var dbobj = new Shop
            {
                Caption = shop.Caption,
                Description = shop.Description,
                DeliveryConditions = shop.DeliveryConditions
            };

            db.Shops.Add(dbobj);
            await db.SaveChangesAsync();

            return CreatedAtAction("GetShop", new { id = dbobj.ID }, dbobj);
        }

        [HttpDelete("{id}")]
        public async Task<IActionResult> DeleteShop(long id)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var shop = await db.Shops.FirstOrDefaultAsync(_ => _.IsDeleted == false && _.ID == id);
            if (shop == null) return NotFound();
            if (shop.CreatedByID != cu.Id) return Unauthorized();

            var clStatuses = new[] { OrderStatuses.Canceled, OrderStatuses.Done };

#pragma warning disable CS8602 // Dereference of a possibly null reference.
            var orders = db.OrderLines
                .Where(_ => _.Shop == shop && _.Order != null)
                .Select(_ => _.Order.ID);
#pragma warning restore CS8602 // Dereference of a possibly null reference.

            var foo = db.Orders.Where(_ => orders.Contains(_.ID) && !_.IsDeleted && !clStatuses.Contains(_.Status)).FirstOrDefault();
            if (foo != null) return BadRequest("Перед удалением витрины нужно закрыть все заказы по ней");

            var goods = db.Goods.Where(_ => _.OwnerShop == shop && !_.IsDeleted);

            using (var scope = sp.CreateScope())
            {
                var gc = new GoodsController(scope.ServiceProvider.GetRequiredService<ApplicationDbContext>(), he, blobServiceClient, cache);
                foreach (var g in goods)
                {
                    var t = gc.DeleteGood(g.ID);
                    t.Wait();
                }
            }

            shop.IsDeleted = true;
            await db.SaveChangesAsync();

            return NoContent();
        }
    }
}
