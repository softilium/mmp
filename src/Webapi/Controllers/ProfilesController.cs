using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.DbCtx;
using mmp.Models;

namespace Webapi.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class ProfilesController(ApplicationDbContext context) : ControllerBase
    {
        private readonly ApplicationDbContext db = context;

        [HttpGet("my")]
        [Authorize]
        public async Task<ActionResult<User>> GetMyProfile()
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var dbuser = await db.Users.FirstOrDefaultAsync(_ => _.Id == cu.Id);
            if (dbuser == null) return NotFound();

            return dbuser;
        }

        [HttpGet("{id}")]
        [Authorize]
        public async Task<ActionResult<UserInfo>> GetUser(long id)
        {
            var user = await db.Users.FindAsync(id);
            if (user == null) return NotFound();
            return UserCache.FindUserInfo(id, db);
        }

        [HttpGet("public")]
        public async Task<ActionResult<UserInfo>> GetPublicUser([FromQuery] string email)
        {
            var user = db.Users
                .Where(_ => _.Email == email)
                .Select(_ => new UserInfo(_))
                .FirstOrDefault();
            if (user == null) return NotFound();
            return user;
        }

        [HttpPut]
        [Authorize]
        public async Task<IActionResult> PutUser(User user)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            // establish uniquess of username
            var newUserName = user.UserName;
            if (string.IsNullOrWhiteSpace(newUserName)) newUserName = user.Email;
            if (db.Users.Any(_ => _.UserName == newUserName && _.Id != cu.Id)) return BadRequest("non-unique user name");

            cu.PhoneNumber = user.PhoneNumber;
            cu.UserName = newUserName;
            cu.Email = user.Email;
            await db.SaveChangesAsync();

            return NoContent();
        }
    }
}
