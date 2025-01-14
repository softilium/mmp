using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.Data;

namespace Webapi.Controllers
{
    [Route("api/admin")]
    [ApiController]
    public class AdminController : ControllerBase
    {
        private readonly ApplicationDbContext db;

        public AdminController(ApplicationDbContext context)
        {
            db = context;
        }

        [HttpGet("allusers")]
        public async Task<ActionResult<IEnumerable<mmp.Data.User>>> GetUsers()
        {
            var cu = db.CurrentUser();
            if (cu == null || !cu.Admin) return Unauthorized();

            return await db.Users.ToListAsync();
        }

        [HttpPut("allusers/{id}")]
        public async Task<IActionResult> PutUser(long id, mmp.Data.User user)
        {
            var cu = db.CurrentUser();
            if (cu == null || !cu.Admin) return Unauthorized();

            if (id != user.Id) return BadRequest();

            var dbUser = db.Users.FirstOrDefault(_ => _.Id == id);
            if (dbUser == null) return NotFound();

            dbUser.Admin = user.Admin;
            dbUser.ShopManage = user.ShopManage;
            await db.SaveChangesAsync();

            return NoContent();
        }

    }
}
