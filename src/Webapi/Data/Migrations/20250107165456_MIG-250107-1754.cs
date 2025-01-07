using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace WebApi.Data.Migrations
{
    /// <inheritdoc />
    public partial class MIG2501071754 : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Goods_AspNetUsers_CreatedById",
                table: "Goods");

            migrationBuilder.DropForeignKey(
                name: "FK_Goods_AspNetUsers_DeletedById",
                table: "Goods");

            migrationBuilder.DropForeignKey(
                name: "FK_Goods_AspNetUsers_ModifiedById",
                table: "Goods");

            migrationBuilder.DropForeignKey(
                name: "FK_OrderLines_AspNetUsers_CreatedById",
                table: "OrderLines");

            migrationBuilder.DropForeignKey(
                name: "FK_OrderLines_AspNetUsers_DeletedById",
                table: "OrderLines");

            migrationBuilder.DropForeignKey(
                name: "FK_OrderLines_AspNetUsers_ModifiedById",
                table: "OrderLines");

            migrationBuilder.DropForeignKey(
                name: "FK_Orders_AspNetUsers_CreatedById",
                table: "Orders");

            migrationBuilder.DropForeignKey(
                name: "FK_Orders_AspNetUsers_DeletedById",
                table: "Orders");

            migrationBuilder.DropForeignKey(
                name: "FK_Orders_AspNetUsers_ModifiedById",
                table: "Orders");

            migrationBuilder.DropForeignKey(
                name: "FK_Shops_AspNetUsers_CreatedById",
                table: "Shops");

            migrationBuilder.DropForeignKey(
                name: "FK_Shops_AspNetUsers_DeletedById",
                table: "Shops");

            migrationBuilder.DropForeignKey(
                name: "FK_Shops_AspNetUsers_ModifiedById",
                table: "Shops");

            migrationBuilder.DropIndex(
                name: "IX_Shops_DeletedById",
                table: "Shops");

            migrationBuilder.DropIndex(
                name: "IX_Shops_ModifiedById",
                table: "Shops");

            migrationBuilder.DropIndex(
                name: "IX_Orders_DeletedById",
                table: "Orders");

            migrationBuilder.DropIndex(
                name: "IX_Orders_ModifiedById",
                table: "Orders");

            migrationBuilder.DropIndex(
                name: "IX_OrderLines_DeletedById",
                table: "OrderLines");

            migrationBuilder.DropIndex(
                name: "IX_OrderLines_ModifiedById",
                table: "OrderLines");

            migrationBuilder.DropIndex(
                name: "IX_Goods_DeletedById",
                table: "Goods");

            migrationBuilder.DropIndex(
                name: "IX_Goods_ModifiedById",
                table: "Goods");

            migrationBuilder.RenameColumn(
                name: "ModifiedById",
                table: "Shops",
                newName: "ModifiedByID");

            migrationBuilder.RenameColumn(
                name: "DeletedById",
                table: "Shops",
                newName: "DeletedByID");

            migrationBuilder.RenameColumn(
                name: "CreatedById",
                table: "Shops",
                newName: "CreatedByID");

            migrationBuilder.RenameIndex(
                name: "IX_Shops_CreatedById",
                table: "Shops",
                newName: "IX_Shops_CreatedByID");

            migrationBuilder.RenameColumn(
                name: "ModifiedById",
                table: "Orders",
                newName: "ModifiedByID");

            migrationBuilder.RenameColumn(
                name: "DeletedById",
                table: "Orders",
                newName: "DeletedByID");

            migrationBuilder.RenameColumn(
                name: "CreatedById",
                table: "Orders",
                newName: "CreatedByID");

            migrationBuilder.RenameIndex(
                name: "IX_Orders_CreatedById",
                table: "Orders",
                newName: "IX_Orders_CreatedByID");

            migrationBuilder.RenameColumn(
                name: "ModifiedById",
                table: "OrderLines",
                newName: "ModifiedByID");

            migrationBuilder.RenameColumn(
                name: "DeletedById",
                table: "OrderLines",
                newName: "DeletedByID");

            migrationBuilder.RenameColumn(
                name: "CreatedById",
                table: "OrderLines",
                newName: "CreatedByID");

            migrationBuilder.RenameIndex(
                name: "IX_OrderLines_CreatedById",
                table: "OrderLines",
                newName: "IX_OrderLines_CreatedByID");

            migrationBuilder.RenameColumn(
                name: "ModifiedById",
                table: "Goods",
                newName: "ModifiedByID");

            migrationBuilder.RenameColumn(
                name: "DeletedById",
                table: "Goods",
                newName: "DeletedByID");

            migrationBuilder.RenameColumn(
                name: "CreatedById",
                table: "Goods",
                newName: "CreatedByID");

            migrationBuilder.RenameIndex(
                name: "IX_Goods_CreatedById",
                table: "Goods",
                newName: "IX_Goods_CreatedByID");

            migrationBuilder.AddColumn<bool>(
                name: "ShopManage",
                table: "AspNetUsers",
                type: "bit",
                nullable: false,
                defaultValue: false);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "ShopManage",
                table: "AspNetUsers");

            migrationBuilder.RenameColumn(
                name: "ModifiedByID",
                table: "Shops",
                newName: "ModifiedById");

            migrationBuilder.RenameColumn(
                name: "DeletedByID",
                table: "Shops",
                newName: "DeletedById");

            migrationBuilder.RenameColumn(
                name: "CreatedByID",
                table: "Shops",
                newName: "CreatedById");

            migrationBuilder.RenameIndex(
                name: "IX_Shops_CreatedByID",
                table: "Shops",
                newName: "IX_Shops_CreatedById");

            migrationBuilder.RenameColumn(
                name: "ModifiedByID",
                table: "Orders",
                newName: "ModifiedById");

            migrationBuilder.RenameColumn(
                name: "DeletedByID",
                table: "Orders",
                newName: "DeletedById");

            migrationBuilder.RenameColumn(
                name: "CreatedByID",
                table: "Orders",
                newName: "CreatedById");

            migrationBuilder.RenameIndex(
                name: "IX_Orders_CreatedByID",
                table: "Orders",
                newName: "IX_Orders_CreatedById");

            migrationBuilder.RenameColumn(
                name: "ModifiedByID",
                table: "OrderLines",
                newName: "ModifiedById");

            migrationBuilder.RenameColumn(
                name: "DeletedByID",
                table: "OrderLines",
                newName: "DeletedById");

            migrationBuilder.RenameColumn(
                name: "CreatedByID",
                table: "OrderLines",
                newName: "CreatedById");

            migrationBuilder.RenameIndex(
                name: "IX_OrderLines_CreatedByID",
                table: "OrderLines",
                newName: "IX_OrderLines_CreatedById");

            migrationBuilder.RenameColumn(
                name: "ModifiedByID",
                table: "Goods",
                newName: "ModifiedById");

            migrationBuilder.RenameColumn(
                name: "DeletedByID",
                table: "Goods",
                newName: "DeletedById");

            migrationBuilder.RenameColumn(
                name: "CreatedByID",
                table: "Goods",
                newName: "CreatedById");

            migrationBuilder.RenameIndex(
                name: "IX_Goods_CreatedByID",
                table: "Goods",
                newName: "IX_Goods_CreatedById");

            migrationBuilder.CreateIndex(
                name: "IX_Shops_DeletedById",
                table: "Shops",
                column: "DeletedById");

            migrationBuilder.CreateIndex(
                name: "IX_Shops_ModifiedById",
                table: "Shops",
                column: "ModifiedById");

            migrationBuilder.CreateIndex(
                name: "IX_Orders_DeletedById",
                table: "Orders",
                column: "DeletedById");

            migrationBuilder.CreateIndex(
                name: "IX_Orders_ModifiedById",
                table: "Orders",
                column: "ModifiedById");

            migrationBuilder.CreateIndex(
                name: "IX_OrderLines_DeletedById",
                table: "OrderLines",
                column: "DeletedById");

            migrationBuilder.CreateIndex(
                name: "IX_OrderLines_ModifiedById",
                table: "OrderLines",
                column: "ModifiedById");

            migrationBuilder.CreateIndex(
                name: "IX_Goods_DeletedById",
                table: "Goods",
                column: "DeletedById");

            migrationBuilder.CreateIndex(
                name: "IX_Goods_ModifiedById",
                table: "Goods",
                column: "ModifiedById");

            migrationBuilder.AddForeignKey(
                name: "FK_Goods_AspNetUsers_CreatedById",
                table: "Goods",
                column: "CreatedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);

            migrationBuilder.AddForeignKey(
                name: "FK_Goods_AspNetUsers_DeletedById",
                table: "Goods",
                column: "DeletedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Goods_AspNetUsers_ModifiedById",
                table: "Goods",
                column: "ModifiedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);

            migrationBuilder.AddForeignKey(
                name: "FK_OrderLines_AspNetUsers_CreatedById",
                table: "OrderLines",
                column: "CreatedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);

            migrationBuilder.AddForeignKey(
                name: "FK_OrderLines_AspNetUsers_DeletedById",
                table: "OrderLines",
                column: "DeletedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_OrderLines_AspNetUsers_ModifiedById",
                table: "OrderLines",
                column: "ModifiedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);

            migrationBuilder.AddForeignKey(
                name: "FK_Orders_AspNetUsers_CreatedById",
                table: "Orders",
                column: "CreatedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);

            migrationBuilder.AddForeignKey(
                name: "FK_Orders_AspNetUsers_DeletedById",
                table: "Orders",
                column: "DeletedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Orders_AspNetUsers_ModifiedById",
                table: "Orders",
                column: "ModifiedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);

            migrationBuilder.AddForeignKey(
                name: "FK_Shops_AspNetUsers_CreatedById",
                table: "Shops",
                column: "CreatedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);

            migrationBuilder.AddForeignKey(
                name: "FK_Shops_AspNetUsers_DeletedById",
                table: "Shops",
                column: "DeletedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Shops_AspNetUsers_ModifiedById",
                table: "Shops",
                column: "ModifiedById",
                principalTable: "AspNetUsers",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);
        }
    }
}
