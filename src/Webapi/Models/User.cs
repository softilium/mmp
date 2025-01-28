using Microsoft.AspNetCore.Identity;
using System.ComponentModel.DataAnnotations.Schema;
using System.ComponentModel.DataAnnotations;
using Microsoft.EntityFrameworkCore;
using System.Text.Json.Serialization;

namespace mmp.Data
{
    // User's projection for frontend
    public class UserInfo
    {
        public string UserName { get; set; }
        public bool ShopManage { get; set; }
        public bool Admin { get; set; }
        public bool TelegramVerified { get; set; }
        public long BotChatId { get; set; }
        public long Id { get; set; }

        public string Description { get; set; } = ""; // by def we don't put value here except public profiles page

        public UserInfo(User src)
        {
            ArgumentNullException.ThrowIfNull(src);
            UserName = src.UserName ?? "";
            ShopManage = src.ShopManage;
            Admin = src.Admin;
            Id = src.Id;
            TelegramVerified = src.TelegramVerified;
            BotChatId = src.BotChatId;
        }

        [JsonConstructor] 
        public UserInfo(string userName, bool shopManage, bool admin, bool telegramVerified, long botChatId, long id) 
        {
            UserName = userName;
            ShopManage = shopManage;
            Admin = admin;
            TelegramVerified = telegramVerified;
            BotChatId = botChatId;
            Id = id;
        }
    }

    public class User : IdentityUser<long>
    {
        public bool ShopManage { get; set; }

        public bool Admin { get; set; }

        [MaxLength(50)]
        public string TelegramUserName { get; set; } = "";

        [MaxLength(20)]

        public string TelegramCheckCode { get; set; } = "";

        public bool TelegramVerified { get; set; } = false;

        [NotMapped]
        public long BotChatId { get; set; } = 0; // shows according chatId from BotChats (for profile page)

        [MaxLength(300)]
        public string Description { get; set; } = "";

    }

    public static class UserCache
    {

        private static Lock lck = new();

        private static Dictionary<long, UserInfo> loaded = [];

        public static void Clear()
        {
            lock (lck) loaded.Clear();
        }
        public static UserInfo FindUserInfo(long id, ApplicationDbContext db)
        {
            lock (lck)
            {
                if (loaded == null || loaded.Count == 0 || !loaded.ContainsKey(id))
                {
                    var chats = db.BotChats.AsNoTracking().ToDictionary(k => k.UserName, v => v.ChatId);
                    var users2tg = db.Users.AsNoTracking().ToDictionary(k => k.UserName, v => v.TelegramUserName); //map user names to tg names
                    loaded = db.Users.AsNoTracking().ToDictionary(k => k.Id, v => new UserInfo(v));
                    foreach (var kv in loaded)
                    {
                        if (chats.TryGetValue(users2tg[kv.Value.UserName], out long chatId))
                        kv.Value.BotChatId = chatId;
                    }
                }

                if (!loaded.TryGetValue(id, out UserInfo? value))
                    throw new Exception($"Unable to find user with id={id}");

                return value;
            }
        }

        public static void LoadCreatedBy(BaseObject obj, ApplicationDbContext db)
        {
            obj.CreatedByInfo = FindUserInfo(obj.CreatedByID, db);
        }
    }

}
