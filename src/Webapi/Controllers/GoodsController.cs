using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.DbCtx;
using mmp.Models;

namespace Webapi.Controllers
{
    [Route("api/goods")]
    [ApiController]
    public class GoodsController : ControllerBase
    {
        private readonly ApplicationDbContext db;

        public GoodsController(ApplicationDbContext _db)
        {
            db = _db;
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<Good>>> GetGoods([FromQuery] long shopId)
        {
            return await db.Goods.Where(_ => _.OwnerShop.ID == shopId).AsNoTracking().ToListAsync();
        }

        [HttpGet("{id}")]
        public async Task<ActionResult<Good>> GetGood(long id)
        {
            var good = await db.Goods.AsNoTracking().FirstOrDefaultAsync(_ => _.ID == id);
            if (good == null) return NotFound();
            return good;
        }

        [HttpPut("{id}")]
        public async Task<IActionResult> PutGood(long id, Good good)
        {
            if (id != good.ID) return BadRequest();
            var dbGood = await db.Goods.FirstOrDefaultAsync(_=>_.ID == id && !_.IsDeleted);
            if (dbGood == null) return NotFound();

            dbGood.Caption = good.Caption;
            dbGood.Description = good.Description;

            await db.SaveChangesAsync();

            return NoContent();
        }

        [HttpPost]
        public async Task<ActionResult<Good>> PostGood(Good good)
        {
            var cu = db.CurrentUser(); 
            if (cu == null) return Unauthorized();

            var shop = await db.Shops.FirstOrDefaultAsync(_=>_.ID==good.OwnerShop.ID && !_.IsDeleted); 
            if (shop == null) return NotFound();

            if (shop.CreatedBy.Email != cu.Email) return Unauthorized();

            var dbGood = new Good
            {
                OwnerShop = shop,
                Caption = good.Caption,
                Description = good.Description
            };

            db.Goods.Add(dbGood);
            await db.SaveChangesAsync();

            return CreatedAtAction("GetGood", new { id = dbGood.ID }, dbGood);
        }

        [HttpDelete("{id}")]
        public async Task<IActionResult> DeleteGood(long id)
        {
            var good = await db.Goods.FirstOrDefaultAsync(_ => _.ID == id && !_.IsDeleted);
            if (good == null) return NotFound();

            good.IsDeleted = true;
            await db.SaveChangesAsync();

            return NoContent();
        }
    }
}
