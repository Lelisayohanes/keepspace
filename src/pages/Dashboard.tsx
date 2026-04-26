import { useState, useEffect } from "react";
import { useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { 
  Plus, 
  Cloud, 
  LogOut, 
  Trash2, 
  Key, 
  ExternalLink,
  Search,
  LayoutDashboard,
  Settings,
  Shield,
  FileText
} from "lucide-react";
import { toast } from "sonner";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { spacesAPI, Space } from "@/lib/api";

export default function Dashboard() {
  const [spaces, setSpaces] = useState<Space[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    fetchSpaces();
  }, []);

  const fetchSpaces = async () => {
    try {
      const data = await spacesAPI.list();
      setSpaces(data);
    } catch (error: any) {
      if (error.message.includes('401') || error.message.includes('Unauthorized')) {
        toast.error("Session expired. Please login again.");
        navigate("/login");
      } else {
        toast.error("Failed to load spaces");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user");
    window.dispatchEvent(new Event("storage"));
    navigate("/login");
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to delete this space? All files will be permanently deleted.")) {
      return;
    }

    try {
      await spacesAPI.delete(id);
      setSpaces(spaces.filter(s => s.id !== id));
      toast.success("Space deleted successfully");
    } catch (error: any) {
      toast.error(error.message || "Failed to delete space");
    }
  };

  return (
    <div className="flex min-h-screen bg-muted/10">
      {/* Sidebar */}
      <aside className="w-64 border-r border-border/50 bg-background hidden md:flex flex-col">
        <div className="h-16 flex items-center px-6 border-b border-border/40">
          <Link className="flex items-center gap-2" to="/dashboard">
            <div className="bg-primary p-1 rounded-lg">
              <Cloud className="h-5 w-5 text-primary-foreground" />
            </div>
            <span className="font-bold text-lg">KeepSpace</span>
          </Link>
        </div>
        <div className="flex-1 py-6 px-4 space-y-2">
          <Button variant="secondary" className="w-full justify-start gap-3 bg-primary/10 text-primary">
            <LayoutDashboard className="h-4 w-4" /> Spaces
          </Button>
          <Button variant="ghost" className="w-full justify-start gap-3">
            <FileText className="h-4 w-4" /> Documentation
          </Button>
          <Button variant="ghost" className="w-full justify-start gap-3">
            <Shield className="h-4 w-4" /> Security
          </Button>
          <Button variant="ghost" className="w-full justify-start gap-3">
            <Settings className="h-4 w-4" /> Settings
          </Button>
        </div>
        <div className="p-4 border-t border-border/40">
          <Button variant="ghost" className="w-full justify-start gap-3 text-destructive hover:text-destructive hover:bg-destructive/10" onClick={handleLogout}>
            <LogOut className="h-4 w-4" /> Logout
          </Button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col">
        <header className="h-16 border-b border-border/40 bg-background flex items-center justify-between px-8 sticky top-0 z-10">
          <h2 className="text-xl font-semibold">Your Spaces</h2>
          <div className="flex items-center gap-4">
            <div className="relative w-64 hidden lg:block">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search spaces..."
                className="pl-9 bg-muted/40 border-none focus-visible:ring-1"
              />
            </div>
            <Button asChild size="sm" className="gap-2 shadow-sm shadow-primary/20">
              <Link to="/spaces/new">
                <Plus className="h-4 w-4" /> Create Space
              </Link>
            </Button>
          </div>
        </header>

        <div className="p-8">
          <div className="grid gap-6 md:grid-cols-3 mb-8">
            <Card className="bg-gradient-to-br from-primary/5 to-transparent border-primary/20">
              <CardHeader className="pb-2">
                <CardDescription>Total Spaces</CardDescription>
                <CardTitle className="text-3xl font-bold">{spaces.length}</CardTitle>
              </CardHeader>
            </Card>
            <Card className="bg-gradient-to-br from-indigo-500/5 to-transparent border-indigo-500/20">
              <CardHeader className="pb-2">
                <CardDescription>Total Files</CardDescription>
                <CardTitle className="text-3xl font-bold">142</CardTitle>
              </CardHeader>
            </Card>
            <Card className="bg-gradient-to-br from-emerald-500/5 to-transparent border-emerald-500/20">
              <CardHeader className="pb-2">
                <CardDescription>Storage Used</CardDescription>
                <CardTitle className="text-3xl font-bold">1.2 GB</CardTitle>
              </CardHeader>
            </Card>
          </div>

          <Card className="border-border/50 shadow-sm overflow-hidden">
            <CardHeader className="bg-muted/30 pb-4">
              <CardTitle className="text-lg">Manage Spaces</CardTitle>
              <CardDescription>View and manage your private storage containers</CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              {isLoading ? (
                <div className="p-8 text-center text-muted-foreground">Loading spaces...</div>
              ) : (
                <Table>
                  <TableHeader className="bg-muted/50">
                    <TableRow>
                      <TableHead className="pl-6">Name</TableHead>
                      <TableHead>ID</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead className="text-right pr-6">Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {spaces.map((space) => (
                      <TableRow key={space.id} className="group hover:bg-muted/20">
                        <TableCell className="font-medium pl-6">
                          <div className="flex items-center gap-2">
                            <Cloud className="h-4 w-4 text-primary/70" />
                            {space.name}
                          </div>
                        </TableCell>
                        <TableCell className="font-mono text-xs text-muted-foreground">{space.id}</TableCell>
                        <TableCell className="text-muted-foreground">{space.created_at}</TableCell>
                        <TableCell>
                          <Badge variant="outline" className="bg-emerald-500/10 text-emerald-600 border-emerald-200">Active</Badge>
                        </TableCell>
                        <TableCell className="text-right pr-6">
                          <div className="flex justify-end gap-2">
                            <Button variant="ghost" size="icon" className="h-8 w-8 hover:text-primary" title="View Keys">
                              <Key className="h-4 w-4" />
                            </Button>
                            <Button variant="ghost" size="icon" className="h-8 w-8 hover:text-primary" title="Open Space">
                              <ExternalLink className="h-4 w-4" />
                            </Button>
                            <Button 
                              variant="ghost" 
                              size="icon" 
                              className="h-8 w-8 text-muted-foreground hover:text-destructive hover:bg-destructive/10"
                              onClick={() => handleDelete(space.id)}
                              title="Delete Space"
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                    {spaces.length === 0 && (
                      <TableRow>
                        <TableCell colSpan={5} className="h-32 text-center text-muted-foreground">
                          No spaces found. Create your first one to get started.
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  );
}