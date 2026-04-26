import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Cloud, Loader2, ArrowLeft, Copy, Check, Key } from "lucide-react";
import { toast } from "sonner";
import { spacesAPI } from "@/lib/api";

export default function CreateSpace() {
  const [name, setName] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [apiKey, setApiKey] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const response = await spacesAPI.create(name);
      setApiKey(response.api_key || null);
      toast.success("Space created successfully!");
    } catch (error: any) {
      toast.error(error.message || "Failed to create space");
      setIsLoading(false);
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = () => {
    if (apiKey) {
      navigator.clipboard.writeText(apiKey);
      setCopied(true);
      toast.success("API Key copied to clipboard");
      setTimeout(() => setCopied(false), 2000);
    }
  };

  if (apiKey) {
    return (
      <div className="min-h-screen flex items-center justify-center px-4 bg-muted/30">
        <Card className="w-full max-w-lg shadow-xl border-emerald-500/20 border-2">
          <CardHeader className="text-center">
            <div className="flex justify-center mb-4">
              <div className="bg-emerald-500 p-2 rounded-xl shadow-inner">
                <Key className="h-8 w-8 text-white" />
              </div>
            </div>
            <CardTitle className="text-2xl font-bold text-emerald-700">Your Space is Ready!</CardTitle>
            <CardDescription className="text-base mt-2">
              Important: Copy your API key now. You will not be able to see it again.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6 pt-4">
            <div className="space-y-3">
              <Label className="text-muted-foreground uppercase text-xs font-bold tracking-widest">Space Name</Label>
              <div className="p-3 bg-muted/50 rounded-lg border border-border/50 font-medium">
                {name}
              </div>
            </div>
            <div className="space-y-3">
              <Label className="text-muted-foreground uppercase text-xs font-bold tracking-widest">API Key</Label>
              <div className="relative group">
                <div className="p-4 bg-background rounded-lg border-2 border-primary/20 font-mono text-sm break-all pr-12 shadow-inner">
                  {apiKey}
                </div>
                <Button 
                  size="icon" 
                  variant="ghost" 
                  className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8"
                  onClick={copyToClipboard}
                >
                  {copied ? <Check className="h-4 w-4 text-emerald-500" /> : <Copy className="h-4 w-4" />}
                </Button>
              </div>
            </div>
            <div className="bg-amber-50 border border-amber-200 p-4 rounded-lg flex gap-3 text-sm text-amber-800">
              <Check className="h-5 w-5 shrink-0 rotate-180" />
              <p>Store this key securely. It provides full programmatic access to this specific Space.</p>
            </div>
          </CardContent>
          <CardFooter>
            <Button className="w-full" onClick={() => navigate("/dashboard")}>
              Done, go to Dashboard
            </Button>
          </CardFooter>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center px-4 bg-muted/30">
      <Card className="w-full max-w-md shadow-lg border-border/50">
        <CardHeader>
          <div className="mb-4">
            <Button variant="ghost" size="sm" asChild className="pl-0 hover:bg-transparent -ml-2">
              <Link to="/dashboard" className="flex items-center gap-1 text-muted-foreground hover:text-primary">
                <ArrowLeft className="h-4 w-4" /> Back to Dashboard
              </Link>
            </Button>
          </div>
          <CardTitle className="text-2xl font-bold">Create New Space</CardTitle>
          <CardDescription>
            Enter a name for your private storage container.
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Space Name</Label>
              <Input
                id="name"
                placeholder="e.g. production-assets"
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="bg-background h-11"
              />
            </div>
            <p className="text-xs text-muted-foreground">
              Spaces are private by default. You'll get an API key to access this Space after creation.
            </p>
          </CardContent>
          <CardFooter>
            <Button className="w-full h-11" type="submit" disabled={isLoading}>
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Create Space
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}