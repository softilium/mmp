using Microsoft.AspNetCore.Identity;
using System.ComponentModel.DataAnnotations.Schema;
using System.ComponentModel.DataAnnotations;

namespace mmp.Data
{
    // User's projection for frontend
    public class UserInfo
    {
        public string UserName { get; set; }
        public bool ShopManage { get; set; }
        public bool Admin { get; set; }
        public bool TelegramVerified { get; set; }    

        public long Id { get; set; }

        public UserInfo(User src)
        {
            ArgumentNullException.ThrowIfNull(src);
            UserName = src.UserName ?? "";
            ShopManage = src.ShopManage;
            Admin = src.Admin;
            Id = src.Id;
            TelegramVerified = src.TelegramVerified;
        }
    }

    public class User : IdentityUser<long>
    {
        public bool ShopManage { get; set; }

        public bool Admin { get; set; }

        [MaxLength(50)]
        public string TelegramUserName { get; set; }

        [MaxLength(20)]

        public string TelegramCheckCode { get; set; }

        public bool TelegramVerified { get; set; }

        [NotMapped]
        public long BotChatId { get; set; } // shows according chatId from BotChats (for profile page)
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
                    loaded = db.Users.ToDictionary(k => k.Id, v => new UserInfo(v));

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
