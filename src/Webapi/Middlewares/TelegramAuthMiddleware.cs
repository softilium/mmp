using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using mmp.Data;
using System.Security.Claims;
using System.Security.Cryptography;
using System.Text;
using System.Web;
using Telegram.Bot;

public class TelegramAuthMiddleWare
{

    private readonly RequestDelegate _next;

    public TelegramAuthMiddleWare(RequestDelegate next)
    {
        _next = next;
    }

    // Telegram users can leave username empty. In this case we use their id with prefix;
    private static string tgUserNameInternal(string UserName, string UserId) => string.IsNullOrWhiteSpace(UserName) ? ("tg." + UserId) : UserName;
    public static string tgUserName(Telegram.Bot.Types.User u) => tgUserNameInternal(u.Username ?? "", u.Id.ToString());

    public static bool CheckInitData(string initData, string botToken)
    {
        try
        {

            // Parse string initData from telegram.
            var data = HttpUtility.ParseQueryString(initData);

            // Put data in a alphabetically sorted dict.
            var dataDict = new SortedDictionary<string, string>(
                data.AllKeys.ToDictionary(x => x!, x => data[x]!),
                StringComparer.Ordinal);

            // Constant key to genrate secret key.
            var constantKey = "WebAppData";

            // https://core.telegram.org/bots/webapps#validating-data-received-via-the-web-app:
            // Data-check-string is a chain of all received fields,
            // sorted alphabetically.
            // in the format key=<value>.
            // with a line feed character ('\n', 0x0A) used as separator.
            // e.g., 'auth_date=<auth_date>\nquery_id=<query_id>\nuser=<user>'
            var dataCheckString = string.Join(
                '\n', dataDict.Where(x => x.Key != "hash") // Hash should be removed.
                    .Select(x => $"{x.Key}={x.Value}")); // like auth_date=<auth_date> ..

            // secrecKey is the HMAC-SHA-256 signature of the bot's token
            // with the constant string WebAppData used as a key.
            var secretKey = HMACSHA256.HashData(
                Encoding.UTF8.GetBytes(constantKey), // WebAppData
                Encoding.UTF8.GetBytes(botToken)); // Bot's token

            var generatedHash = HMACSHA256.HashData(
                secretKey,
                Encoding.UTF8.GetBytes(dataCheckString)); // data_check_string

            // Convert received hash from telegram to a byte array.
            var actualHash = Convert.FromHexString(dataDict["hash"]); // .NET 5.0 

            // Compare our hash with the one from telegram.
            return actualHash.SequenceEqual(generatedHash);
        }
        catch
        {
            Console.WriteLine("Invalig tg token to validate: " + initData);
            return false;
        }
    }

    public async Task Invoke(HttpContext context, TelegramBotClient bot, ApplicationDbContext db, UserManager<User> um)
    {
        if (context.Request.Headers.ContainsKey("Authorization") && context.Request.Headers["Authorization"].Count > 0)
        {
            var tgauth = context.Request.Headers["Authorization"][0] ?? "     ";
            if (tgauth.StartsWith("tg "))
            {
                tgauth = tgauth.Substring(3);

                var tgauthtokens = tgauth.Split("~~", StringSplitOptions.None);
                if (tgauthtokens.Length != 3)
                {
                    Console.WriteLine("Unable to parse tg auth token: " + tgauth);
                    return;
                }
                tgauth = tgauthtokens[0];
                var tgAuthId = tgauthtokens[1];
                var tgAuthUserName = tgauthtokens[2];

                if (!string.IsNullOrWhiteSpace(tgauth))
                {
                    var botToken = Environment.GetEnvironmentVariable("TelegramBotAPIKEY") ?? "nobothere";
                    if (CheckInitData(tgauth, botToken))
                    {
                        var username = tgUserNameInternal(tgAuthUserName, tgAuthId);

                        if (username != "")
                        {
                            var user = await db.Users.AsNoTracking().FirstOrDefaultAsync(_ => _.TelegramUserName == username && _.TelegramVerified);
                            if (user == null)
                            {
                                // lets create default user record for telegram user 

                                var chat = await db.BotChats.FirstOrDefaultAsync(_ => _.UserName == username)
                                    ?? throw new Exception($"Unable to find stored chat for tg username={username}. Does bot active? Does bot handler?");

                                var newUser = new User
                                {
                                    TelegramUserName = username,
                                    UserName = username,
                                    TelegramVerified = true,
                                    Email = username + "@telegram.tg"
                                };

                                var res = await um.CreateAsync(newUser);
                                if (!res.Succeeded)
                                {
                                    throw new Exception(res.Errors.ToString());
                                }

                                var newPassword = Guid.NewGuid().ToString()[^8..];

                                await um.AddPasswordAsync(newUser, newPassword);

                                await db.SaveChangesAsync();

                                await bot.SendMessage(chat.ChatId,

@$"
Вас приветствует бот сервиса RiverStores. Вы можете использовать сервис, нажав на название бота сверху и открыв мини-приложение по ссылке. Либо, в списке чатов или в окне чата нажать кнопку (ОТКРЫТЬ/OPEN).

Вы можете открыть ссылку в браузере: https://rives-stores.com
Ваш логин  : {newUser.Email}
Ваш пароль : {newPassword}

Есть вопросы, проблемы, предложения, идеи, отзывы? Просто напишите их в чат боту и они будут сразу переданы администратору сервиса."
                                );
                                user = newUser;
                            }
                            var identity = new ClaimsIdentity(new[] { new Claim(ClaimTypes.Name, user.UserName ?? "") }, "custom");
                            context.User = new ClaimsPrincipal(identity);
                        }
                    }
                }
            }
        }
        await _next(context);
    }
}
