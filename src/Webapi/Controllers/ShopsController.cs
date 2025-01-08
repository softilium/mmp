using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.DbCtx;
using mmp.Models;

namespace Webapi.Controllers
{
    [Route("api/shops")]
    [ApiController]
    public class ShopsController : ControllerBase
    {
        private readonly ApplicationDbContext db;

        public ShopsController(ApplicationDbContext context)
        {
            db = context;
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
        [Authorize]
        public async Task<IActionResult> PutShop(long id, Shop shop)
        {

            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var dbobj = db.Shops.First(_ => _.ID == id && _.IsDeleted == false && shopManagers().Contains(_.CreatedByID));
            if (dbobj == null) return NotFound();

            if (dbobj.CreatedByID != cu.Id) return Unauthorized();

            dbobj.Caption = shop.Caption;

            await db.SaveChangesAsync();

            return NoContent();
        }

        [HttpPost]
        [Authorize]
        public async Task<ActionResult<Shop>> PostShop(Shop shop)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            if (!shopManagers().Contains(cu.Id)) return Unauthorized();

            var dbobj = new Shop { Caption = shop.Caption };

            db.Shops.Add(dbobj);
            await db.SaveChangesAsync();

            return CreatedAtAction("GetShop", new { id = dbobj.ID }, dbobj);
        }

        [HttpDelete("{id}")]
        [Authorize]
        public async Task<IActionResult> DeleteShop(long id)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var shop = await db.Shops.FirstOrDefaultAsync(_ => _.IsDeleted == false && _.ID == id);
            if (shop == null) return NotFound();
            if (shop.CreatedByID != cu.Id) return Unauthorized();

            shop.IsDeleted = true;
            await db.SaveChangesAsync();

            return NoContent();
        }
    }
}
