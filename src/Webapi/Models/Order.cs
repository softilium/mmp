using Microsoft.EntityFrameworkCore;
using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using Webapi.Controllers;

namespace mmp.Data
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

        [MaxLength(200)]
        public string? SenderComment { get; set; }
        
        [MaxLength(200)]
        public string? CustomerComment { get; set; }

        public DateTime? ExpectedDeliveryDate { get; set; }

        public ICollection<OrderLine> Lines { get; set; } = [];

        public override void BeforeSave(ApplicationDbContext db, Microsoft.EntityFrameworkCore.ChangeTracking.EntityEntry entity)
        {
            if (entity.State != EntityState.Modified) return;

            var oldStatus = entity.OriginalValues[nameof(Status)] == null
                ? Status
                : (OrderStatuses)entity.OriginalValues[nameof(Status)];

            if (Status != oldStatus)
            {
                var clientUser = UserCache.FindUserInfo(CreatedByID, db);
                if (!clientUser.TelegramVerified)
                {
                    var senderUser = UserCache.FindUserInfo(Shop.CreatedByID, db);
                    db.NotifyAfterSave(senderUser.BotChatId, $"Заказчик {clientUser.UserName} для заказа {ID} от {CreatedOn:g} не получает уведомления, не настроена интеграция с Телеграм");
                    return;
                }
                else
                {
                    if (Status != oldStatus)
                        db.NotifyAfterSave(clientUser.BotChatId, $"Статус вашего заказа {ID} от {CreatedOn:g} изменился с [{oldStatus.GetEnumDescription()}] на [{Status.GetEnumDescription()}]");
                }
            }
        }
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
