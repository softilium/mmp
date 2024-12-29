using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;
using mmp.Models;

namespace mmp.DbCtx
{
    public class ApplicationDbContext : IdentityDbContext<User, IdentityRole<long>, long>
    {
        public DbSet<Shop> Shops { get; set; }
        public DbSet<Good> Goods { get; set; }
        public DbSet<Order> Orders { get; set; }
        public DbSet<OrderLine> OrderLines { get; set; }

        private IHttpContextAccessor _ctx;

        public ApplicationDbContext(
            DbContextOptions<ApplicationDbContext> options,
            IHttpContextAccessor ctx
        )
            : base(options)
        {
            _ctx = ctx;

            SavingChanges += (s, e) =>
            {
                if (_ctx.HttpContext == null || _ctx.HttpContext.User.Identity == null || !_ctx.HttpContext.User.Identity.IsAuthenticated)
                    return;
                var currentUser = Users.First(_ =>
                    _.UserName == _ctx.HttpContext.User.Identity.Name
                );

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
                            baseObj.CreatedOn = DateTime.UtcNow;
                            baseObj.CreatedByID = currentUser.Id;
                        }
                        else if (q.State == EntityState.Modified)
                        {
                            baseObj.ModifiedOn = DateTime.UtcNow;
                            baseObj.ModifiedByID = currentUser.Id;
                        }
                        if (
                            baseObj.IsDeleted
                            && (baseObj.DeletedOn == null || baseObj.DeletedByID == null)
                        )
                        {
                            baseObj.DeletedOn = DateTime.UtcNow;
                            baseObj.DeletedByID = currentUser.Id;
                        }
                    }
                }
            };
        }

        protected override void OnModelCreating(ModelBuilder mb)
        {
            base.OnModelCreating(mb);

            //mb.Entity<Models.Quote>().HasKey(_ => new { _.TickerId, _.D });
            //mb.Entity<Models.Quote>()
            //    .HasOne(_ => _.Ticker)
            //    .WithMany(_ => _.Quotes)
            //    .HasForeignKey(_ => _.TickerId);

            //mb.Entity<Models.DSI>().HasKey(_ => new { _.TickerId, _.TargetYear });
            //mb.Entity<Models.DSI>()
            //    .HasOne(_ => _.Ticker)
            //    .WithMany(_ => _.DSI)
            //    .HasForeignKey(_ => _.TickerId);

            //mb.Entity<Models.DivPayout>()
            //    .HasOne(d => d.Ticker)
            //    .WithMany(t => t.DivPayouts)
            //    .HasForeignKey(_ => _.TickerId);
            //mb.Entity<Models.DivPayout>()
            //    .HasKey(_ => new
            //    {
            //        _.TickerId,
            //        _.ForYear,
            //        _.ForQuarter,
            //        _.CloseDate,
            //        _.DPS
            //    });

            //mb.Entity<Models.StrategyFactor>().HasKey(_ => new { _.StrategyId, _.LineNum });

            //mb.Entity<Models.StrategyFilter>().HasKey(_ => new { _.StrategyId, _.LineNum });

            //mb.Entity<Models.StrategyPreselected>().HasKey(_ => new { _.StrategyId, _.LineNum });

            //mb.Entity<Models.Strategy>()
            //    .Property(_ => _.SameEmitent)
            //    .HasDefaultValue(Models.SameEmitentPolicy.Allow);

            //mb.Entity<Models.InvestmentAccount>()
            //    .HasOne(_ => _.Owner)
            //    .WithMany(_ => _.InvestAccounts)
            //    .HasForeignKey(_ => _.OwnerId);

            //mb.Entity<Models.InvestmentAccountEvaluation>()
            //    .HasKey(_ => new { _.AccountId, _.RecDate });
            //mb.Entity<Models.InvestmentAccountEvaluation>()
            //    .HasOne(_ => _.Account)
            //    .WithMany(_ => _.Evaluations)
            //    .HasForeignKey(_ => _.AccountId);

            //mb.Entity<Models.InvestmentAccountCashflow>()
            //    .HasKey(_ => new { _.AccountId, _.RecDate });
            //mb.Entity<Models.InvestmentAccountCashflow>()
            //    .HasOne(_ => _.Account)
            //    .WithMany(_ => _.Cashflow)
            //    .HasForeignKey(_ => _.AccountId);
        }
    }
}
