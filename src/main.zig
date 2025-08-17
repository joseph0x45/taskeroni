const zqlite = @import("zqlite");
const std = @import("std");

const migrationCreate =
    \\create table if not exists users (
    \\id text not null primary key,
    \\username text not null unique,
    \\password text not null,
    \\is_admin integer not null default 0
    \\);
    \\
;

const migrationSeed =
    \\insert or ignore into users (id, username, password, is_admin)
    \\values('1', 'admin', 'admin', 1);
;

const User = struct {
    ID: []const u8,
    Username: []const u8,
    Password: []const u8,
    IsAdmin: bool,
};

pub fn main() !void {
    const flags = zqlite.OpenFlags.Create | zqlite.OpenFlags.EXResCode;
    const connection = try zqlite.open("test.sqlite", flags);
    defer connection.close();

    try connection.exec(migrationCreate, .{});
    try connection.exec(migrationSeed, .{});

    if (try connection.row("select * from users where username=?1", .{"admin"})) |row| {
        defer row.deinit();
        std.debug.print("id: {s}\n", .{row.text(0)});
        std.debug.print("username: {s}\n", .{row.text(1)});
        std.debug.print("password: {s}\n", .{row.text(2)});
        std.debug.print("is_admin: {}\n", .{row.int(3)});
    }
}
