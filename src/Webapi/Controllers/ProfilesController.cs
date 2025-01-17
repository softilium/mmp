using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using mmp.Data;
using System.Runtime.InteropServices;
using Telegram.Bot;

namespace Webapi.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class ProfilesController(ApplicationDbContext context, TelegramBotClient _bot) : ControllerBase
    {
        private readonly ApplicationDbContext db = context;
        private readonly TelegramBotClient bot = _bot;
    
        [HttpGet("my")]
        [Authorize]
        public async Task<ActionResult<User>> GetMyProfile()
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var dbuser = await db.Users.AsNoTracking().FirstOrDefaultAsync(_ => _.Id == cu.Id);
            if (dbuser == null) return NotFound();

            var chat = db.BotChats.AsNoTracking().FirstOrDefault(_ => _.UserName == dbuser.TelegramUserName);
            if (chat != null)
                dbuser.BotChatId = chat.ChatId;

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
            if (string.IsNullOrWhiteSpace(email))
            {
                var cu = db.CurrentUser();
                if (cu == null) return NotFound();
                return new UserInfo(cu);
            }
            var user = await db.Users
                .Where(_ => _.Email == email)
                .Select(_ => new UserInfo(_))
                .FirstOrDefaultAsync();
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

            cu.UserName = newUserName;
            cu.Email = user.Email;
            cu.TelegramUserName = user.TelegramUserName;
            cu.TelegramVerified = user.TelegramVerified;
            cu.TelegramCheckCode = user.TelegramCheckCode;

            await db.SaveChangesAsync();

            return NoContent();
        }

        [HttpPost("newtelegramcode")]
        [Authorize]
        public async Task<IActionResult> SetNewCodeForTelegram()
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            var chat = db.BotChats.FirstOrDefault(_ => _.UserName == cu.UserName);
            if (chat == null) return NotFound();

            var nc = Guid.NewGuid().ToString();
            nc = nc.Substring(nc.Length - 5, 5).ToUpper();
            cu.TelegramCheckCode = nc;

            await db.SaveChangesAsync();

            await bot.SendMessage(chat.ChatId,
                $"{nc}\n\r\n\rЭто ваш проверочный код для подтверждения в учетной записи на сайте river-stores.com. Он действует один раз, не передавайте его никому."
            );

            return Ok();
        }

        [HttpPost("checktelegramcode/{newcode}")]
        [Authorize]
        public async Task<IActionResult> CheckCodeForTelegram(string newcode)
        {
            var cu = db.CurrentUser();
            if (cu == null) return Unauthorized();

            if (!cu.TelegramCheckCode.Trim().Equals(newcode.Trim(), StringComparison.CurrentCultureIgnoreCase)) return BadRequest();

            var chat = db.BotChats.FirstOrDefault(_ => _.UserName == cu.UserName);
            if (chat == null) return NotFound();

            cu.TelegramVerified = true;
            cu.TelegramCheckCode = "";

            await db.SaveChangesAsync();

            await bot.SendMessage(chat.ChatId,
                $"Ваш проверочный код подтвержден на сайте и этот чат будет в дальнейшем использоваться для уведомлений о событиях на сайте river-store.com для вас"
            );

            return Ok();
        }

    }
}
