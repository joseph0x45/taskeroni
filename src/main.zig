const zqlite = @import("zqlite");
const std = @import("std");
const httpz = @import("httpz");
const buildin = @import("builtin");

const indexHtml = @embedFile("static/index.html");
const authHtml = @embedFile("static/auth.html");
const stylesCss = @embedFile("static/styles.css");
const scriptJs = @embedFile("static/script.js");
const alpineJs = @embedFile("static/alpine.js");

const PORT = 8080;

const User = struct {
    id: []const u8,
    username: []const u8,
};

const App = struct {
    db: *const zqlite.Conn,
};

const File = struct {
    content: []const u8,
    mimeType: []const u8,
};

const migrationCreate =
    \\create table if not exists users (
    \\id integer not null primary key,
    \\username text not null unique
    \\);
;

const migrationSeed =
    \\insert or ignore into users (id, username)
    \\values(1, 'default');
;

pub fn main() !void {
    const flags = zqlite.OpenFlags.Create | zqlite.OpenFlags.EXResCode;
    const connection = try zqlite.open("test.sqlite", flags);
    defer connection.close();

    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    try connection.exec(migrationCreate, .{});
    try connection.exec(migrationSeed, .{});

    var app: App = App{
        .db = &connection,
    };
    var server = try httpz.Server(*App).init(allocator, .{ .port = PORT }, &app);
    var router = try server.router(.{});

    router.get("/", renderHomePage, .{});
    router.get("/auth", renderAuthPage, .{});
    router.get("/static/:file", serveStaticFiles, .{});

    std.debug.print("Server Listening on port {}\nVisit http://0.0.0.0:8080/auth\n", .{PORT});
    try server.listen();
}

fn serveStaticFiles(_: *App, req: *httpz.Request, res: *httpz.Response) !void {
    if (req.param("file")) |fileName| {
        res.status = 200;
        if (std.mem.eql(u8, fileName, "styles.css")) {
            res.headers.add("Content-Type", "text/css");
            try res.writer().writeAll(stylesCss);
            return;
        } else if (std.mem.eql(u8, fileName, "alpine.js")) {
            res.headers.add("Content-Type", "application/javascript");
            try res.writer().writeAll(alpineJs);
            return;
        } else if (std.mem.eql(u8, fileName, "script.js")) {
            res.headers.add("Content-Type", "application/javascript");
            try res.writer().writeAll(scriptJs);
            return;
        }
    }
    res.status = 404;
    try res.writer().print("Not Found", .{});
}

fn renderHomePage(_: *App, _: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    res.headers.add("Content-Type", "text/html");
    try res.writer().writeAll(indexHtml);
}

fn renderAuthPage(_: *App, _: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    res.headers.add("Content-Type", "text/html");
    try res.writer().writeAll(authHtml);
}
