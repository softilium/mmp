using Microsoft.EntityFrameworkCore;
using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
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
        public long SenderID { get; set; }

        [NotMapped]
        public UserInfo? SenderInfo { get; set; }

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

        public string DeltaTxt(string ov, string nv) =>
            string.IsNullOrWhiteSpace(ov)
            ? $"\n\r\n\r{nv}"
            : $"\n\r\n\rБЫЛО:\n\r{ov}\n\r\n\rСТАЛО:\n\r{nv}";

        public override void BeforeSave(ApplicationDbContext db, Microsoft.EntityFrameworkCore.ChangeTracking.EntityEntry entity)
        {
            if (entity.State == EntityState.Added)
            {
                var cu = db.CurrentUser();

                var senderUser = UserCache.FindUserInfo(SenderID, db);

                db.NotifyAfterSave(senderUser.BotChatId, $"Новый заказ для вас. Заказчик {cu.UserName}, {DateTime.Now:g}.");

                if (cu.TelegramVerified)
                    db.NotifyAfterSave(cu.BotChatId, $"Ваш заказ #{ID} от {CreatedOn: g} отправлен владельцу витрины {senderUser.UserName}. Бот уведомит вас обо всех изменениях его статуса.");

            }

            if (entity.State != EntityState.Modified) return;

            var oldStatus = (OrderStatuses)entity.OriginalValues[nameof(Status)];

            var oldExpectedDeliveryDate = entity.OriginalValues[nameof(ExpectedDeliveryDate)] == null
                ? null
                : (DateTime?)entity.OriginalValues[nameof(ExpectedDeliveryDate)];

            var oldCustomerComment = (string)entity.OriginalValues[nameof(CustomerComment)];
            var oldSenderComment = (string)entity.OriginalValues[nameof(SenderComment)];

            if (Status != oldStatus || oldExpectedDeliveryDate != ExpectedDeliveryDate || oldCustomerComment != CustomerComment || oldSenderComment != SenderComment)
            {
                var senderUser = UserCache.FindUserInfo(SenderID, db);
                var clientUser = UserCache.FindUserInfo(CreatedByID, db);
                var cu = db.CurrentUser();
                if (cu == null) return;

                if (cu.Id == CreatedByID)
                {
                    if (oldCustomerComment != CustomerComment)
                        db.NotifyAfterSave(senderUser.BotChatId, $"Заказчик {clientUser.UserName} по заказу {ID} от {CreatedOn:g} указал примечание к заказу.{DeltaTxt(oldCustomerComment ?? "", CustomerComment ?? "")}");
                }
                if (cu.Id == SenderID)
                {
                    if (!clientUser.TelegramVerified)
                    {
                        db.NotifyAfterSave(senderUser.BotChatId, $"Заказчик {clientUser.UserName} для заказа {ID} от {CreatedOn:g} не получает уведомления, не настроена интеграция с Телеграм");
                        return;
                    }
                    if (Status != oldStatus)
                        db.NotifyAfterSave(clientUser.BotChatId, $"Статус вашего заказа {ID} от {CreatedOn:g} изменился.{DeltaTxt(oldStatus.GetEnumDescription(), Status.GetEnumDescription())}");

                    if (ExpectedDeliveryDate != oldExpectedDeliveryDate)
                        db.NotifyAfterSave(clientUser.BotChatId, $"Ожидаемая дата доставки заказа {ID} от {CreatedOn:g} уточнена.{DeltaTxt(
                            oldExpectedDeliveryDate == null ? "": oldExpectedDeliveryDate.Value.ToString("g"),
                            ExpectedDeliveryDate == null ? "" : ExpectedDeliveryDate.Value.ToString("g")
                        )}");

                    if (oldSenderComment != SenderComment)
                        db.NotifyAfterSave(clientUser.BotChatId, $"Отправитель {senderUser.UserName} по заказу {ID} от {CreatedOn:g} указал примечание продавца.{DeltaTxt(oldSenderComment ?? "", SenderComment ?? "")}");
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

