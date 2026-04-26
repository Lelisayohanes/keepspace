import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import { Shield, Cloud, Key, ArrowRight, Github } from "lucide-react";
import { motion } from "framer-motion";

export default function Landing() {
  return (
    <div className="flex flex-col min-h-screen">
      <header className="px-4 lg:px-6 h-16 flex items-center border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky top-0 z-50">
        <Link className="flex items-center justify-center gap-2" to="/">
          <div className="bg-primary p-1.5 rounded-lg">
            <Cloud className="h-6 w-6 text-primary-foreground" />
          </div>
          <span className="text-xl font-bold tracking-tight">KeepSpace</span>
        </Link>
        <nav className="ml-auto flex gap-4 sm:gap-6 items-center">
          <Link className="text-sm font-medium hover:text-primary transition-colors" to="/login">
            Login
          </Link>
          <Button asChild size="sm">
            <Link to="/signup">Get Started</Link>
          </Button>
        </nav>
      </header>
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32 xl:py-48">
          <div className="container px-4 md:px-6 mx-auto">
            <div className="grid gap-6 lg:grid-cols-[1fr_400px] lg:gap-12 xl:grid-cols-[1fr_600px] items-center">
              <motion.div 
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5 }}
                className="flex flex-col justify-center space-y-4"
              >
                <div className="space-y-2">
                  <h1 className="text-3xl font-bold tracking-tighter sm:text-5xl xl:text-6xl/none">
                    Your Data, Your Control. <span className="text-primary">KeepSpace</span>.
                  </h1>
                  <p className="max-w-[600px] text-muted-foreground md:text-xl">
                    The minimal, secure, self-hosted S3-compatible storage gateway. Manage your files with private Spaces and programmatic API keys.
                  </p>
                </div>
                <div className="flex flex-col gap-2 min-[400px]:flex-row">
                  <Button asChild size="lg" className="px-8">
                    <Link to="/signup">
                      Deploy Your Space <ArrowRight className="ml-2 h-4 w-4" />
                    </Link>
                  </Button>
                  <Button variant="outline" size="lg" className="px-8">
                    <Github className="mr-2 h-4 w-4" /> View Source
                  </Button>
                </div>
              </motion.div>
              <motion.div
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.5, delay: 0.2 }}
              >
                <img
                  alt="KeepSpace Dashboard Preview"
                  className="mx-auto aspect-video overflow-hidden rounded-xl object-cover object-center shadow-2xl sm:w-full lg:order-last border border-border"
                  src="https://storage.googleapis.com/dala-prod-public-storage/generated-images/7147aa4f-7147-437c-b5ae-29ea02e8cbd1/dashboard-preview-6e4c76fd-1777228188793.webp"
                />
              </motion.div>
            </div>
          </div>
        </section>
        <section className="w-full py-12 md:py-24 lg:py-32 bg-muted/50">
          <div className="container px-4 md:px-6 mx-auto">
            <div className="flex flex-col items-center justify-center space-y-4 text-center">
              <div className="space-y-2">
                <div className="inline-block rounded-lg bg-primary/10 px-3 py-1 text-sm text-primary">Features</div>
                <h2 className="text-3xl font-bold tracking-tighter sm:text-5xl">Built for Developers</h2>
                <p className="max-w-[900px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                  Simple, powerful, and private. Everything you need to manage object storage without the complexity.
                </p>
              </div>
            </div>
            <div className="mx-auto grid max-w-5xl items-center gap-6 py-12 lg:grid-cols-3 lg:gap-12">
              <div className="flex flex-col items-center space-y-4 text-center">
                <div className="bg-background p-4 rounded-full shadow-sm border border-border">
                  <Shield className="h-10 w-10 text-primary" />
                </div>
                <h3 className="text-xl font-bold">Private Spaces</h3>
                <p className="text-muted-foreground">Isolate your data in secure containers. No buckets, just Spaces.</p>
              </div>
              <div className="flex flex-col items-center space-y-4 text-center">
                <div className="bg-background p-4 rounded-full shadow-sm border border-border">
                  <Key className="h-10 w-10 text-primary" />
                </div>
                <h3 className="text-xl font-bold">API Key Access</h3>
                <p className="text-muted-foreground">Generate programmatic keys for each Space. Secure and easy to rotate.</p>
              </div>
              <div className="flex flex-col items-center space-y-4 text-center">
                <div className="bg-background p-4 rounded-full shadow-sm border border-border">
                  <Cloud className="h-10 w-10 text-primary" />
                </div>
                <h3 className="text-xl font-bold">S3 Powered</h3>
                <p className="text-muted-foreground">Built on top of MinIO/S3. High performance and industry standard compatibility.</p>
              </div>
            </div>
          </div>
        </section>
      </main>
      <footer className="flex flex-col gap-2 sm:flex-row py-6 w-full shrink-0 items-center px-4 md:px-6 border-t border-border/40">
        <p className="text-xs text-muted-foreground">© 2024 KeepSpace Inc. All rights reserved.</p>
        <nav className="sm:ml-auto flex gap-4 sm:gap-6">
          <Link className="text-xs hover:underline underline-offset-4 text-muted-foreground" to="#">
            Terms of Service
          </Link>
          <Link className="text-xs hover:underline underline-offset-4 text-muted-foreground" to="#">
            Privacy
          </Link>
        </nav>
      </footer>
    </div>
  );
}