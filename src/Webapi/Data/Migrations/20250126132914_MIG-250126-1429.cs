using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace WebApi.Data.Migrations
{
    /// <inheritdoc />
    public partial class MIG2501261429 : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Orders_Shops_ShopID",
                table: "Orders");

            migrationBuilder.DropIndex(
                name: "IX_Orders_ShopID",
                table: "Orders");

            migrationBuilder.RenameColumn(
                name: "ShopID",
                table: "Orders",
                newName: "SenderID");

            migrationBuilder.Sql("delete from public.\"OrderLines\"");
            migrationBuilder.Sql("delete from public.\"Orders\"");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.RenameColumn(
                name: "SenderID",
                table: "Orders",
                newName: "ShopID");

            migrationBuilder.CreateIndex(
                name: "IX_Orders_ShopID",
                table: "Orders",
                column: "ShopID");

            migrationBuilder.AddForeignKey(
                name: "FK_Orders_Shops_ShopID",
                table: "Orders",
                column: "ShopID",
                principalTable: "Shops",
                principalColumn: "ID",
                onDelete: ReferentialAction.Restrict);
        }
    }
}
