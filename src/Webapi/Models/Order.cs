using Microsoft.EntityFrameworkCore;
using System.ComponentModel;
using System.ComponentModel.DataAnnotations;

namespace mmp.Models
{
    public enum OrderStatuses : int
    {
        [Description("Новый")]
        New = 100,

        [Description("В работе")]
        InProcess = 200,

        [Description("Готов к доставке")]
        ReadyToDeliver = 300,

        [Description("Доставляется")]
        Delivering = 400,

        [Description("Доставлен")]
        Done = 500,

        [Description("Отменен")]
        Canceled = 600,
    }

    public class Order : BaseObject
    {
        [Required]
        [DeleteBehavior(DeleteBehavior.Restrict)] 
        public Shop Shop { get; set; }

        [Required] 
        public OrderStatuses Status { get; set; }

        [Required]
        [Precision(15, 2)] 
        public decimal Qty { get; set; }

        [Required]
        [Precision(15, 2)] 
        public decimal Sum { get; set; }

        public ICollection<OrderLine> Lines { get; set; } = [];
    }

    public class OrderLine : BaseObject
    {
        [Required]
        [DeleteBehavior(DeleteBehavior.Restrict)] 
        public Shop Shop { get; set; }
        
        public Order? Order { get; set; } //when empty it is basket
        
        [Required]
        [DeleteBehavior(DeleteBehavior.Restrict)] 
        public Good Good { get; set; }
        
        [Required]
        [Precision(15, 2)] 
        public decimal Qty { get; set; }
        
        [Required]
        [Precision(15, 2)] 
        public decimal Sum { get; set; }
    }
}
