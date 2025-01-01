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
        private readonly ApplicationDbContext ctx;
        private readonly IHttpContextAccessor req;

        public ShopsController(ApplicationDbContext context, IHttpContextAccessor RequestCtx)
        {
            ctx = context;
            req = RequestCtx;
        }

        // GET: api/Shops
        [HttpGet]
        public async Task<ActionResult<IEnumerable<Shop>>> GetShops()
        {
            return await ctx.Shops.Where(_ => _.IsDeleted == false).ToListAsync();
        }

        // GET: api/Shops/5
        [HttpGet("{id}")]
        public async Task<ActionResult<Shop>> GetShop(long id)
        {
            var shop = await ctx.Shops
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

            var dbobj = ctx.Shops.First(_ => _.ID == id && _.IsDeleted == false);
            if (dbobj == null) return NotFound();

            dbobj.OwnerUserID = shop.OwnerUserID;
            dbobj.Name = shop.Name;
            dbobj.Caption = shop.Caption;
            dbobj.Comment = shop.Comment;

            await ctx.SaveChangesAsync();

            return NoContent();
        }

        // POST: api/Shops
        // To protect from overposting attacks, see https://go.microsoft.com/fwlink/?linkid=2123754
        [HttpPost]
        public async Task<ActionResult<Shop>> PostShop(Shop shop)
        {

            var IsAuth = 
                req.HttpContext != null && req.HttpContext.User.Identity != null && req.HttpContext.User.Identity.IsAuthenticated;
            if (!IsAuth) return Unauthorized();

            var dbobj = new Shop
            {
                OwnerUserID = shop.OwnerUserID,
                Name = shop.Name,
                Caption = shop.Caption
            };

            ctx.Shops.Add(dbobj);
            await ctx.SaveChangesAsync();

            return CreatedAtAction("GetShop", new { id = dbobj.ID }, dbobj);
        }

        // DELETE: api/Shops/5
        [HttpDelete("{id}")]
        public async Task<IActionResult> DeleteShop(long id)
        {
            var shop = await ctx.Shops.FirstOrDefaultAsync(_=>_.IsDeleted==false && _.ID==id);
            if (shop == null)
            {
                return NotFound();
            }

            ctx.Shops.Remove(shop);
            await ctx.SaveChangesAsync();

            return NoContent();
        }
    }
}
