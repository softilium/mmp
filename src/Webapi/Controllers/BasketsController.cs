using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.DbCtx;
using mmp.Models;
using System.Net.NetworkInformation;

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
        [Route("{shopId:long}")]
        [Authorize]
        public async Task<ActionResult<IEnumerable<OrderLine>>> GetOrderLines(long shopId)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            return await db.OrderLines
                .Include(_ => _.Good)
                .Where(_ => _.CreatedBy.Id == cu.Id && _.Shop.ID == shopId && _.Order == null)
                .AsNoTracking()
                .ToListAsync();
        }

        [HttpPost]
        [Authorize]
        [Route("increase/{goodId:long}")]
        public async Task<ActionResult<OrderLine>> Increase(long goodId, [FromQuery] decimal? qty)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var good = db.Goods.Include(_ => _.OwnerShop).FirstOrDefault(_ => _.ID == goodId);
            if (good == null) return NotFound();

            qty = qty ?? 1;

            var ol = db.OrderLines.Where(_ => _.CreatedBy.Id == cu.Id && _.Order == null && _.Good.ID == goodId).FirstOrDefault();
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
        [Authorize]
        [Route("decrease/{goodId:long}")]
        public async Task<ActionResult<OrderLine>> Decrease(long goodId, [FromQuery] decimal? qty)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            qty = qty ?? 1;

            var ol = db.OrderLines.Where(_ => _.CreatedBy.Id == cu.Id && _.Order == null && _.Good.ID == goodId).FirstOrDefault();
            if (ol != null)
            {
                var price = ol.Sum / ol.Qty;
                ol.Qty -= qty.Value;
                ol.Sum -= ol.Qty * price;
                if (ol.Qty <= 0) db.OrderLines.Remove(ol);
                await db.SaveChangesAsync();
                return NoContent();
            }
            return NoContent();
        }

    }
}
