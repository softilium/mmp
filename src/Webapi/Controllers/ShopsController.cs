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
        private readonly IHttpContextAccessor req;

        public ShopsController(ApplicationDbContext context, IHttpContextAccessor RequestCtx)
        {
            db = context;
            req = RequestCtx;
        }

        // GET: api/Shops
        [HttpGet]
        public async Task<ActionResult<IEnumerable<Shop>>> GetShops()
        {
            return await db.Shops.Where(_ => _.IsDeleted == false).ToListAsync();
        }

        // GET: api/Shops/5
        [HttpGet("{id}")]
        public async Task<ActionResult<Shop>> GetShop(long id)
        {
            var shop = await db.Shops
                .Where(_ => _.IsDeleted == false && _.ID == id)
                .FirstOrDefaultAsync();

            if (shop == null) return NotFound();

            return shop;
        }

        // PUT: api/Shops/5
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPut("{id}")]
        public async Task<IActionResult> PutShop(long id, Shop shop)
        {

            var dbobj = db.Shops.First(_ => _.ID == id && _.IsDeleted == false);
            if (dbobj == null) return NotFound();

            dbobj.Caption = shop.Caption;
            dbobj.Comment = shop.Comment;

            await db.SaveChangesAsync();

            return NoContent();
        }

        [HttpPost]
        [Authorize]
        public async Task<ActionResult<Shop>> PostShop(Shop shop)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var dbobj = new Shop
            {
                Caption = shop.Caption
            };

            db.Shops.Add(dbobj);
            await db.SaveChangesAsync();

            return CreatedAtAction("GetShop", new { id = dbobj.ID }, dbobj);
        }

        // DELETE: api/Shops/5
        [HttpDelete("{id}")]
        public async Task<IActionResult> DeleteShop(long id)
        {
            var shop = await db.Shops.FirstOrDefaultAsync(_ => _.IsDeleted == false && _.ID == id);
            if (shop == null)
            {
                return NotFound();
            }

            db.Shops.Remove(shop);
            await db.SaveChangesAsync();

            return NoContent();
        }
    }
}
