const zqlite = @import("zqlite");
const std = @import("std");
const httpz = @import("httpz");

const PORT = 8080;

const indexHTML: []const u8 = @embedFile("static/index.html");
const stylesCSS: []const u8 = @embedFile("static/styles.css");
const scriptJS: []const u8 = @embedFile("static/script.js");

const User = struct {
    id: []const u8,
    username: []const u8,
    password: []const u8,
    isAdmin: bool,
};

const App = struct {
    db: *const zqlite.Conn,
};

const File = struct {
    path: []const u8,
    content: []const u8,
    mimeType: []const u8,
};

const embeddedFiles: [3]File = .{
    File{ .path = "index.html", .content = indexHTML, .mimeType = "text/html" },
    File{ .path = "styles.css", .content = stylesCSS, .mimeType = "text/css" },
    File{ .path = "script.js", .content = scriptJS, .mimeType = "applicatin/javascript" },
};

const migrationCreate =
    \\create table if not exists users (
    \\id text not null primary key,
    \\username text not null unique,
    \\password text not null,
    \\is_admin integer not null default 0
    \\);
;

const migrationSeed =
    \\insert or ignore into users (id, username, password, is_admin)
    \\values('1', 'admin', 'admin', 1);
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

    router.get("/hello", sayHello, .{});
    router.get("/index.html", serveHTML, .{});
    router.get("/static/:file", serveStaticFiles, .{});
    std.debug.print("Server Listening on port {}\n", .{PORT});
    try server.listen();
}

fn sayHello(_: *App, _: *httpz.Request, res: *httpz.Response) !void {
    try res.json(.{ .message = "World" }, .{});
}

fn serveHTML(_: *App, _: *httpz.Request, res: *httpz.Response) !void {
    res.status = 200;
    res.headers.add("Content-Type", "text/html");
    try res.writer().writeAll(indexHTML);
}

fn serveStaticFiles(_: *App, req: *httpz.Request, res: *httpz.Response) !void {
    var found = false;
    const fileName = req.param("file") orelse "index.html";
    for (embeddedFiles) |f| {
        if (std.mem.eql(u8, f.path, fileName)) {
            res.headers.add("Cache-Control", "no-store, no-cache, must-revalidate");
            res.headers.add("Pragma", "no-cache");
            res.headers.add("Expires", "0");
            res.headers.add("Content-Type", f.mimeType);
            res.status = 200;
            try res.writer().writeAll(f.content);
            found = true;
            break;
        }
    }
    if (!found) {
        res.status = 404;
        try res.writer().print("Not Found", .{});
    }
}
