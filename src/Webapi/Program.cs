using Microsoft.EntityFrameworkCore;
using mmp.Data;
using Microsoft.AspNetCore.Authentication;
using System.Text.Json.Serialization;
using Azure.Storage.Blobs;
using Telegram.Bot;
using Telegram.Bot.Types.Enums;
using ZiggyCreatures.Caching.Fusion;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddFusionCache()
    .WithDefaultEntryOptions(opt => { opt.Duration = TimeSpan.FromHours(1); });

builder.Configuration.AddEnvironmentVariables(); // azure uses env.variables for app config

{
    var connectionString = builder.Configuration.GetConnectionString("DefaultConnection");
    if (string.IsNullOrWhiteSpace(connectionString))
        connectionString = builder.Configuration["DefaultConnection"];

    if (string.IsNullOrWhiteSpace(connectionString))
        Console.WriteLine("ERROR. Unable to get DefaultConnection");

    builder.Services.AddDbContext<ApplicationDbContext>(options =>
    {
        options.UseNpgsql(connectionString);
        options.EnableSensitiveDataLogging(builder.Environment.IsDevelopment());
    });
}

{
    var storageAccountConnStr = builder.Configuration.GetConnectionString("StorageAccount");
    if (string.IsNullOrWhiteSpace(storageAccountConnStr))
        storageAccountConnStr = builder.Configuration["StorageAccount"];
    if (string.IsNullOrWhiteSpace(storageAccountConnStr))
        Console.WriteLine("ERROR. Enable to get StorageAccount");

    builder.Services.AddSingleton(x => new BlobServiceClient(storageAccountConnStr));
}

AppContext.SetSwitch("Npgsql.EnableLegacyTimestampBehavior", true);

var TelegramBotAPIKEY = builder.Configuration["TelegramBotAPIKEY"] ?? "";

builder.Services.AddSingleton(x => new TelegramBotClient(TelegramBotAPIKEY));

//builder.Services.AddAuthorization();

builder.Services.AddIdentityApiEndpoints<User>()
    .AddEntityFrameworkStores<ApplicationDbContext>();

builder.Services.AddOpenApi();

builder.Services.AddCors(options =>
{
    options.AddPolicy(name: "MyPolicy",
        b => { b.WithOrigins("*").AllowAnyMethod().AllowAnyOrigin().AllowAnyHeader(); });
});

builder.Services.AddControllers()
    .AddJsonOptions(options => options.JsonSerializerOptions.ReferenceHandler = ReferenceHandler.IgnoreCycles);

var app = builder.Build();

app.UseMiddleware<TelegramAuthMiddleWare>();

app.UseCors("MyPolicy");

app.MapGroup("identity").MapIdentityApi<User>();
app.MapPost("/identity/logout", ctx => ctx.SignOutAsync()); // std. MapIdentityApi doesn't contains logout endpoint

if (app.Environment.IsDevelopment())
    app.MapOpenApi();
else
    app.UseHttpsRedirection();

app.MapControllers();

using (var scope = app.Services.CreateScope())
{
    var db = scope.ServiceProvider.GetRequiredService<ApplicationDbContext>();
    db.Database.Migrate();
}

#region bot

if (!string.IsNullOrEmpty(TelegramBotAPIKEY))
{
    var bot = new TelegramBotClient(TelegramBotAPIKEY);
    var me = await bot.GetMe();

    if (!app.Environment.IsDevelopment())
        bot.OnMessage += OnMessage;

    // associate users with chatIds
    async Task OnMessage(Telegram.Bot.Types.Message msg, UpdateType type)
    {
        if (msg == null || msg.From == null) return;

        using var botscope = app.Services.CreateScope();

        var db = botscope.ServiceProvider.GetRequiredService<ApplicationDbContext>();
        var chat = db.BotChats.FirstOrDefault(_ => _.ChatId == msg.Chat.Id);
        if (chat == null)
        {
            chat = new BotChat { ChatId = msg.Chat.Id };
            db.BotChats.Add(chat);
            chat.UserName = TelegramAuthMiddleWare.tgUserName(msg.From);
            await db.SaveChangesAsync();
            Console.WriteLine($"New chat {msg.Chat.Id} is associated with {TelegramAuthMiddleWare.tgUserName(msg.From)}");
        }
        else
        {
            if (chat.UserName != TelegramAuthMiddleWare.tgUserName(msg.From))
            {
                chat.UserName = TelegramAuthMiddleWare.tgUserName(msg.From);
                await db.SaveChangesAsync();
                Console.WriteLine($"Updated chat {msg.Chat.Id}. Now it's associated with {TelegramAuthMiddleWare.tgUserName(msg.From)}");
            }
        }
        var admUserNames = db.Users.AsNoTracking().Where(_ => _.Admin && _.TelegramVerified).Select(_ => _.TelegramUserName).ToList();
        var admChats = db.BotChats.AsNoTracking().Where(_ => admUserNames.Contains(_.UserName)).Select(_ => _.ChatId).ToList();
        if (!admChats.Contains(msg.Chat.Id))
        {
            db.NotifyFirstAdminAfterSave($"Message from {msg.Chat.Username}:\n\r{msg.Text}");
            await db.SaveChangesAsync();
        }
    }
}
else Console.WriteLine("Bot isn't initialized, TelegramBotAPIKEY is empty");

#endregion

app.Run();
