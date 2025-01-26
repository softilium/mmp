using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.Data;

namespace Webapi.Controllers
{
    [Route("api/baskets")]
    [ApiController]
    public class BasketsController : ControllerBase
    {
        private readonly ApplicationDbContext db;

        public BasketsController(ApplicationDbContext _context)
        {
            db = _context;
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<OrderLine>>> GetWholeBasket()
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();


            var res = await db.OrderLines
                .AsNoTracking()
                .Include(_ => _.Good)
                .Where(_ => _.CreatedByID == cu.Id && _.Order == null)
                .ToListAsync();

            foreach (var line in res)
                UserCache.LoadCreatedBy(line.Good, db);

            return res;
        }

        [HttpGet]
        [Route("{shopId:long}")]
        public async Task<ActionResult<IEnumerable<OrderLine>>> GetBasketForShop(long shopId)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            return await db.OrderLines
                .Include(_ => _.Good)
                .Where(_ => _.CreatedByID == cu.Id && _.Shop.ID == shopId && _.Order == null)
                .AsNoTracking()
                .ToListAsync();
        }

        [HttpPost]
        [Route("increase/{goodId:long}")]
        public async Task<ActionResult<OrderLine>> Increase(long goodId, [FromQuery] decimal? qty)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var good = db.Goods.Include(_ => _.OwnerShop).FirstOrDefault(_ => _.ID == goodId);
            if (good == null) return NotFound();

            qty = qty ?? 1;

            var ol = db.OrderLines
                .Where(_ => _.CreatedByID == cu.Id && _.Order == null && _.Good.ID == goodId)
                .FirstOrDefault();
            if (ol != null)
            {
                ol.Qty += qty.Value;
                ol.Sum += qty.Value * good.Price;
                await db.SaveChangesAsync();
                return NoContent();
            }
            else
            {
                ol = new()
                {
                    Good = good,
                    Shop = good.OwnerShop,
                    Qty = qty.Value,
                    Sum = qty.Value * good.Price
                };
                db.OrderLines.Add(ol);
                await db.SaveChangesAsync();
                return NoContent();
            }
        }

        [HttpPost]
        [Route("decrease/{goodId:long}")]
        public async Task<ActionResult<OrderLine>> Decrease(long goodId, [FromQuery] decimal? qty)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            qty = qty ?? 1;

            var ol = db.OrderLines
                .Include(_=>_.Good)
                .Where(_ => _.CreatedByID == cu.Id && _.Order == null && _.Good.ID == goodId)
                .FirstOrDefault();
            if (ol != null)
            {
                ol.Qty -= qty.Value;
                ol.Sum = ol.Qty * ol.Good.Price;
                if (ol.Qty <= 0) db.OrderLines.Remove(ol);
                await db.SaveChangesAsync();
                return NoContent();
            }
            return NoContent();
        }

    }
}
