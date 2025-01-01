using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using mmp.DbCtx;
using mmp.Models;
using Microsoft.AspNetCore.Authentication;

var builder = WebApplication.CreateBuilder(args);

var connectionString = builder.Configuration.GetConnectionString("DefaultConnection");
builder.Services.AddDbContext<ApplicationDbContext>(options =>
{
    options.UseSqlServer(connectionString);
    options.EnableSensitiveDataLogging(builder.Environment.IsDevelopment());
});

builder.Services.AddIdentityApiEndpoints<User>(opt =>
{
    opt.Password.RequiredLength = 1; //todo : change password min.length to 8
    opt.User.RequireUniqueEmail = true;
    opt.Password.RequireNonAlphanumeric = false;
    opt.SignIn.RequireConfirmedEmail = false;
})
        .AddEntityFrameworkStores<ApplicationDbContext>()
        .AddRoles<IdentityRole<long>>()
    .AddEntityFrameworkStores<ApplicationDbContext>();

builder.Services.AddAuthentication();

builder.Services.AddAuthorization();

builder.Services.AddOpenApi();

builder.Services.AddCors(options =>
{
    options.AddPolicy(name: "MyPolicy",
        b =>
        {
            //This is how you tell your app to allow cors
            b
                .WithOrigins("*")
                //.WithMethods("POST", "DELETE", "GET", "OPTIONS")
                .AllowAnyMethod()
                .AllowAnyOrigin()
                .AllowAnyHeader();
        });
});

builder.Services.AddControllers();

var app = builder.Build();

app.UseCors("MyPolicy");
app.UseAuthentication();
app.UseAuthorization();


app.MapGroup("identity").MapIdentityApi<User>();
app.MapPost("/identity/logout", ctx => ctx.SignOutAsync()); // std. MapIdentityApi doesn't contains logout endpoint

if (app.Environment.IsDevelopment())
    app.MapOpenApi();
else
    app.UseHttpsRedirection();

app.MapControllers();

using (var scope = app.Services.CreateScope())
{
    var db = scope.ServiceProvider.GetRequiredService<ApplicationDbContext>();
    db.Database.Migrate();
}
app.Run();
