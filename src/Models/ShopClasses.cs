using System.ComponentModel.DataAnnotations;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
namespace mmp.Models
{
    public class User : IdentityUser<long> { }

    public abstract class BaseObject
    {
        [Key] public long ID { get; set; }

        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public User CreatedBy { get; set; } = new();
        [Required] public DateTime CreatedOn { get; set; }

        [DeleteBehavior(DeleteBehavior.Restrict)] public User? ModifiedBy { get; set; }
        public DateTime ModifiedOn { get; set; }

        public bool IsDeleted { get; set; } = false;
        public User? DeletedBy { get; set; }
        public DateTime? DeletedOn { get; set; }

        public string? Comment { get; set; }

        public void BeforeSave() { }
    }

    public class Shop : BaseObject
    {
        [Required] public string Caption { get; set; } = "Shop 1";

    }

    public class Good : BaseObject
    {
        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public Shop OwnerShop { get; set; } = new();
        [Required] public string Caption { get; set; } = "";
        public string? Description { get; set; }
        public decimal Price { get; set; }
    }

    public enum OrderStatuses : int
    {
        New = 100,
        ReadyToDeliver = 200,
        Delivered = 300,
        Canceled = 400,
    }

    public class Order : BaseObject
    {
        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public Shop? Shop { get; set; }
        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public User? Buyer { get; set; }
        [Required] public OrderStatuses Status { get; set; }
    }

    public class OrderLine : BaseObject
    {
        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public User Who { get; set; } = new();
        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public Shop Shop { get; set; } = new();
        public Order? Order { get; set; } //when empty it is basket
        [Required][DeleteBehavior(DeleteBehavior.Restrict)] public Good Good { get; set; } = new();
        [Required] public long Qty { get; set; }
    }
}

//todo tags
//todo categories
//todo Prices columns
//todo integration-telegram-payments
//todo comments
//todo search (elastic)
//todo telegram for notify
//todo blacklist
//todo voting for goods/shops (rating, ...)
//todo recommendations
//todo location for shops/delivery
//todo want.to.buy.later (notify me when price go low)
//todo notify all past customers
//todo docker
//todo web blazor client
//todo maui ios
//todo maui android
//todo telegram webapp
//todo azure sql
//todo azure funcs
//todo azure app services
//todo azure static web apps
//todo xunit, nunit
//todo typescript
//todo redis
//todo llm
//todo mini-profiler
//todo opened-closed shops (public, direct-link, approved-only)
//todo dataversion and checking for it before update
