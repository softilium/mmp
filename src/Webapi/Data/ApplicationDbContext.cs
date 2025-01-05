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

        public ApplicationDbContext(
            DbContextOptions<ApplicationDbContext> options,
            IHttpContextAccessor ctx
        )
            : base(options)
        {
            _ctx = ctx;
            SavingChanges += (s, e) =>
            {
                var currentUser = CurrentUser();
                if (currentUser != null)
                {

                    foreach (var q in ChangeTracker.Entries())
                    {
                        if (
                            q.Entity is BaseObject baseObj
                            && (q.State == EntityState.Added || q.State == EntityState.Modified)
                        )
                        {
                            baseObj.BeforeSave();

                            if (q.State == EntityState.Added)
                            {
                                baseObj.CreatedOn = DateTime.Now;
                                baseObj.CreatedBy = currentUser;
                            }
                            else if (q.State == EntityState.Modified)
                            {
                                baseObj.ModifiedOn = DateTime.Now;
                                baseObj.ModifiedBy = currentUser;
                            }
                            if (
                                baseObj.IsDeleted
                                && (baseObj.DeletedOn == null || baseObj.DeletedBy == null)
                            )
                            {
                                baseObj.DeletedOn = DateTime.Now;
                                baseObj.DeletedBy = currentUser;
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
