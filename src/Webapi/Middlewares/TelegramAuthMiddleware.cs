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

    public static bool CheckInitData(string initData, string botToken, out string username)
    {
        //Console.WriteLine($"initData={initData}");
        username = "";

        // Parse string initData from telegram.
        var data = HttpUtility.ParseQueryString(initData);

        // Put data in a alphabetically sorted dict.
        var dataDict = new SortedDictionary<string, string>(
            data.AllKeys.ToDictionary(x => x!, x => data[x]!),
            StringComparer.Ordinal);

        foreach (var kv in dataDict)
        {
            //Console.WriteLine($"{kv.Key} = {kv.Value}");
            if (kv.Key == "user")
            {   // extract value from "username" field in JSON
                var usernamePos = kv.Value.IndexOf("username") + 9;
                var firstQ = kv.Value.IndexOf("\"", usernamePos) + 1;
                var secondQ = kv.Value.IndexOf("\"", firstQ);
                username = kv.Value[firstQ..secondQ];
            }
        }

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

    public async Task Invoke(HttpContext context, TelegramBotClient bot, ApplicationDbContext db)
    {
        if (context.Request.Headers.ContainsKey("tgauth") && context.Request.Headers["tgauth"].Count > 0)
        {
            var tgauth = context.Request.Headers["tgauth"][0];
            if (!string.IsNullOrWhiteSpace(tgauth))
            {
                var botToken = Environment.GetEnvironmentVariable("TelegramBotAPIKEY");
                var username = "";
                if (CheckInitData(tgauth, botToken, out username))
                {
                    var user = await db.Users.FirstOrDefaultAsync(_ => _.TelegramUserName == username);
                    if (user != null)
                    {
                        var claims = new[] { new Claim("name", user.UserName) };
                        var identity = new ClaimsIdentity(claims, "Basic");
                        context.User = new ClaimsPrincipal(identity);
                    }
                }
            }
        }
        await _next(context);
    }
}

