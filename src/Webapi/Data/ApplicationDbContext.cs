using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;
using mmp.Models;

namespace mmp.DbCtx
{
    public class ApplicationDbContext : IdentityDbContext<User, IdentityRole<long>, long>
    {

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

        public ApplicationDbContext(
            DbContextOptions<ApplicationDbContext> options,
            IHttpContextAccessor ctx
        )
            : base(options)
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
                if (currentUser != null)
                {
                    foreach (var q in ChangeTracker.Entries())
                    {
                        if (q.Entity is BaseObject baseObj && (q.State == EntityState.Added || q.State == EntityState.Modified))
                        {
                            baseObj.BeforeSave();

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
        }

        protected override void OnModelCreating(ModelBuilder mb)
        {
            base.OnModelCreating(mb);
        }
    }

}
