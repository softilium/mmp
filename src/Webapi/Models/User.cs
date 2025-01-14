using Microsoft.AspNetCore.Identity;
using System.ComponentModel.DataAnnotations.Schema;
using System.ComponentModel.DataAnnotations;

namespace mmp.Models
{
    // User's projection for frontend
    public class UserInfo
    {
        public string UserName { get; set; }
        public bool ShopManage { get; set; }
        public bool Admin { get; set; }

        public long Id { get; set; }

        public UserInfo(User src)
        {
            ArgumentNullException.ThrowIfNull(src);
            UserName = src.UserName ?? "";
            ShopManage = src.ShopManage;
            Admin = src.Admin;
            Id = src.Id;
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
}
