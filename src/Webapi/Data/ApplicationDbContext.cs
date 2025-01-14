using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;
using Telegram.Bot;

namespace mmp.Data
{
    public class ApplicationDbContext : IdentityDbContext<User, IdentityRole<long>, long>
    {

        private List<long> adminChatIds = null;

        #region NotifyAfterSave
        private readonly Dictionary<long, List<string>> afterSaveNotifies = [];
        public void NotifyAfterSave(long chatId, string message)
        {
            if (!afterSaveNotifies.TryGetValue(chatId, out List<string>? value))
            {
                value = [];
                afterSaveNotifies.Add(chatId, value);
            }
            value.Add(message);
        }

        public void NotifyAdminsAfterSave(string message)
        {
            if (adminChatIds == null)
            {
                var allAdmins = Users.Where(_ => _.Admin && _.TelegramVerified).Select(_ => _.TelegramUserName);
                adminChatIds = BotChats.Where(_ => allAdmins.Contains(_.UserName)).Select(_ => _.ChatId).ToList();
            }
            foreach (var chatId in adminChatIds)
                NotifyAfterSave(chatId, $"АДМ.СООБЩЕНИЕ:\n\r{message}");
        }

        public void ClearAfterSaveNotifies()
        {
            afterSaveNotifies.Clear();
        }
        #endregion

        private IHttpContextAccessor _ctx;
        public User? CurrentUser()
        {
            if (_ctx == null || _ctx.HttpContext == null || _ctx.HttpContext.User == null || _ctx.HttpContext.User.Identity == null) return null;
            if (!_ctx.HttpContext.User.Identity.IsAuthenticated) return null;
            return Users.FirstOrDefault(_ => _.UserName == _ctx.HttpContext.User.Identity.Name);
        }

        public DbSet<Shop> Shops { get; set; }
        public DbSet<Good> Goods { get; set; }
        public DbSet<Order> Orders { get; set; }
        public DbSet<OrderLine> OrderLines { get; set; }
        public DbSet<BotChat> BotChats { get; set; }

        public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options, IHttpContextAccessor ctx, TelegramBotClient _bot) : base(options)
        {
            _ctx = ctx;
            SavingChanges += (s, e) =>
            {
                // invalidate user cache when we save any user
                foreach (var q in ChangeTracker.Entries())
                    if (q.Entity is User && (q.State == EntityState.Added || q.State == EntityState.Modified))
                    {
                        UserCache.Clear();
                        break;
                    }

                // make first user admin
                foreach (var q in ChangeTracker.Entries())
                    if (q.Entity is User userEntity && (q.State == EntityState.Added))
                        if (!Users.Any()) userEntity.Admin = true;
            };

            SavingChanges += (s, e) =>
            {
                var currentUser = CurrentUser();
                foreach (var q in ChangeTracker.Entries())
                {
                    if (q.Entity is BaseObject baseObj && (q.State == EntityState.Added || q.State == EntityState.Modified))
                    {
                        baseObj.BeforeSave(this, q);

                        if (currentUser != null)
                        {
                            if (q.State == EntityState.Added)
                            {
                                baseObj.CreatedOn = DateTime.Now;
                                baseObj.CreatedByID = currentUser.Id;
                            }
                            else if (q.State == EntityState.Modified)
                            {
                                baseObj.ModifiedOn = DateTime.Now;
                                baseObj.ModifiedByID = currentUser.Id;
                            }
                            if (baseObj.IsDeleted && (baseObj.DeletedOn == null || baseObj.DeletedByID == null))
                            {
                                baseObj.DeletedOn = DateTime.Now;
                                baseObj.DeletedByID = currentUser.Id;
                            }
                        }
                    }
                }
            };

            SavedChanges += (s, e) =>
            {
                foreach (var kp in afterSaveNotifies)
                    _bot.SendMessage(kp.Key, string.Join("\n\r", kp.Value.ToArray()));
                ClearAfterSaveNotifies();
            };
        }

        protected override void OnModelCreating(ModelBuilder mb)
        {
            base.OnModelCreating(mb);
        }
    }

}
