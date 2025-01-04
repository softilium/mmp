using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.DbCtx;
using mmp.Models;

namespace Webapi.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class BasketsController : ControllerBase
    {
        private readonly ApplicationDbContext db;

        public BasketsController(ApplicationDbContext _context)
        {
            db = _context;
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<OrderLine>>> GetOrderLines(long shopId)
        {
            return await db.OrderLines.Where(_ => _.Shop.ID == shopId && _.Order == null).AsNoTracking().ToListAsync();
        }

        [HttpPost]
        [Authorize]
        [Route("increase/{goodId:long}")]
        public async Task<ActionResult<OrderLine>> Increase(long goodId, [FromQuery] long? qty)
        {

            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            qty = qty ?? 1;

            var ol = db.OrderLines.Where(_ => _.Order == null && _.Good.ID == goodId).FirstOrDefault();
            if (ol != null)
            {
                ol.Qty += qty.Value;
                await db.SaveChangesAsync();
                return NoContent();
            }
            else
            {
                var good = db.Goods.FirstOrDefault(_ => _.ID == goodId);
                if (good == null) return NotFound();

                ol = new();
                ol.Who.Id = cu.Id;
                ol.Good.ID = goodId;
                ol.Shop.ID = good.OwnerShop.ID;
                db.OrderLines.Add(ol);

                await db.SaveChangesAsync();

                return NoContent();
            }
        }

        [HttpPost]
        [Authorize]
        [Route("decrease/{goodId:long}")]
        public async Task<ActionResult<OrderLine>> Decrease(long goodId, [FromQuery] long? qty)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            qty = qty ?? 1;

            var ol = db.OrderLines.Where(_ => _.Order == null && _.Good.ID == goodId).FirstOrDefault();
            if (ol != null)
            {
                ol.Qty -= qty.Value;
                if (ol.Qty <= 0) db.OrderLines.Remove(ol);
                await db.SaveChangesAsync();
                return NoContent();
            }
            return NoContent();
        }

    }
}
