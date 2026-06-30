package templates

import (
	"fmt"
	"os"
	"path/filepath"
)

// Template represents a supported project template.
type Template struct {
	Name        string
	Description string
}

// GetList returns the list of supported templates.
func GetList() []Template {
	return []Template{
		{Name: "springboot", Description: "Spring Boot starter (Java with Maven pom.xml and main REST controller)"},
		{Name: "react", Description: "React app using Vite (Vanilla JS and basic template structure)"},
		{Name: "rust", Description: "Rust CLI Application with Cargo setup"},
		{Name: "go", Description: "Go API Service with net/http router and handlers folder"},
		{Name: "fastapi", Description: "Python FastAPI application with Dockerfile"},
		{Name: "nextjs", Description: "Next.js App Router (TypeScript boilerplate setup)"},
		{Name: "php", Description: "Composer-based PHP web routing starter skeleton"},
		{Name: "dotnet", Description: ".NET Core C# Web API minimal API setup"},
		{Name: "java", Description: "Standard Java Application with Gradle build config"},
	}
}

// Generate project template by name under target directory.
func Generate(templateName, projectName, destDir string) error {
	projectPath := filepath.Join(destDir, projectName)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	switch templateName {
	case "springboot":
		return generateSpringBoot(projectPath, projectName)
	case "react":
		return generateReact(projectPath, projectName)
	case "rust":
		return generateRust(projectPath, projectName)
	case "go":
		return generateGo(projectPath, projectName)
	case "fastapi":
		return generateFastAPI(projectPath, projectName)
	case "nextjs":
		return generateNextJS(projectPath, projectName)
	case "php":
		return generatePHP(projectPath, projectName)
	case "dotnet":
		return generateDotnet(projectPath, projectName)
	case "java":
		return generateJava(projectPath, projectName)
	default:
		return fmt.Errorf("unsupported template: %s", templateName)
	}
}

func generateSpringBoot(root, name string) error {
	// Write pom.xml
	pom := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
	<modelVersion>4.0.0</modelVersion>
	<parent>
		<groupId>org.springframework.boot</groupId>
		<artifactId>spring-boot-starter-parent</artifactId>
		<version>3.2.4</version>
		<relativePath/>
	</parent>
	<groupId>com.example</groupId>
	<artifactId>%s</artifactId>
	<version>0.0.1-SNAPSHOT</version>
	<name>%s</name>
	<description>Demo project for Spring Boot</description>
	<properties>
		<java.version>17</java.version>
	</properties>
	<dependencies>
		<dependency>
			<groupId>org.springframework.boot</groupId>
			<artifactId>spring-boot-starter-web</artifactId>
		</dependency>
		<dependency>
			<groupId>org.springframework.boot</groupId>
			<artifactId>spring-boot-starter-test</artifactId>
			<scope>test</scope>
		</dependency>
	</dependencies>
	<build>
		<plugins>
			<plugin>
				<groupId>org.springframework.boot</groupId>
				<artifactId>spring-boot-maven-plugin</artifactId>
			</plugin>
		</plugins>
	</build>
</project>
`, name, name)

	// Java App
	javaDir := filepath.Join(root, "src", "main", "java", "com", "example", name)
	if err := os.MkdirAll(javaDir, 0755); err != nil {
		return err
	}

	appClass := `package com.example.` + name + `;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@SpringBootApplication
@RestController
public class Application {

	public static void main(String[] args) {
		SpringApplication.run(Application.class, args);
	}

	@GetMapping("/")
	public String hello() {
		return "Hello, Developer! Welcome to your Spring Boot REST API.";
	}
}
`
	resDir := filepath.Join(root, "src", "main", "resources")
	if err := os.MkdirAll(resDir, 0755); err != nil {
		return err
	}

	_ = os.WriteFile(filepath.Join(root, "pom.xml"), []byte(pom), 0644)
	_ = os.WriteFile(filepath.Join(javaDir, "Application.java"), []byte(appClass), 0644)
	_ = os.WriteFile(filepath.Join(resDir, "application.properties"), []byte("server.port=8080\n"), 0644)
	return nil
}

func generateReact(root, name string) error {
	packageJSON := fmt.Sprintf(`{
  "name": "%s",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.66",
    "@types/react-dom": "^18.2.22",
    "@vitejs/plugin-react": "^4.2.1",
    "vite": "^5.2.0"
  }
}
`, name)

	viteConfig := `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
})
`

	indexHTML := `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>React App</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.jsx"></script>
  </body>
</html>
`

	appJSX := `import React from 'react'

function App() {
  return (
    <div style={{ textAlign: 'center', fontFamily: 'sans-serif', marginTop: '50px' }}>
      <h1>Welcome to React + Vite</h1>
      <p>Configure and run standard frontend terminals easily.</p>
    </div>
  )
}

export default App
`

	mainJSX := `import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
`

	srcDir := filepath.Join(root, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return err
	}

	_ = os.WriteFile(filepath.Join(root, "package.json"), []byte(packageJSON), 0644)
	_ = os.WriteFile(filepath.Join(root, "vite.config.js"), []byte(viteConfig), 0644)
	_ = os.WriteFile(filepath.Join(root, "index.html"), []byte(indexHTML), 0644)
	_ = os.WriteFile(filepath.Join(srcDir, "App.jsx"), []byte(appJSX), 0644)
	_ = os.WriteFile(filepath.Join(srcDir, "main.jsx"), []byte(mainJSX), 0644)
	return nil
}

func generateRust(root, name string) error {
	cargo := fmt.Sprintf(`[package]
name = "%s"
version = "0.1.0"
edition = "2021"

[dependencies]
`, name)

	main := `fn main() {
    println!("Hello, Developer! This is a Rust CLI template created by reshell.");
}
`
	src := filepath.Join(root, "src")
	if err := os.MkdirAll(src, 0755); err != nil {
		return err
	}

	_ = os.WriteFile(filepath.Join(root, "Cargo.toml"), []byte(cargo), 0644)
	_ = os.WriteFile(filepath.Join(src, "main.rs"), []byte(main), 0644)
	return nil
}

func generateGo(root, name string) error {
	goMod := fmt.Sprintf(`module %s

go 1.22
`, name)

	main := `package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, welcome to your Go API boilerplate!")
	})

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
`
	_ = os.WriteFile(filepath.Join(root, "go.mod"), []byte(goMod), 0644)
	_ = os.WriteFile(filepath.Join(root, "main.go"), []byte(main), 0644)
	return nil
}

func generateFastAPI(root, name string) error {
	main := `from fastapi import FastAPI

app = FastAPI(title="FastAPI Service")

@app.get("/")
def read_root():
    return {"message": "Hello World", "template": "FastAPI"}
`
	reqs := `fastapi>=0.110.0
uvicorn>=0.28.0
`

	dockerfile := `FROM python:3.9-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
`

	_ = os.WriteFile(filepath.Join(root, "main.py"), []byte(main), 0644)
	_ = os.WriteFile(filepath.Join(root, "requirements.txt"), []byte(reqs), 0644)
	_ = os.WriteFile(filepath.Join(root, "Dockerfile"), []byte(dockerfile), 0644)
	return nil
}

func generateNextJS(root, name string) error {
	packageJSON := fmt.Sprintf(`{
  "name": "%s",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start"
  },
  "dependencies": {
    "next": "14.1.4",
    "react": "^18",
    "react-dom": "^18"
  },
  "devDependencies": {
    "@types/node": "^20",
    "@types/react": "^18",
    "@types/react-dom": "^18",
    "typescript": "^5"
  }
}
`, name)

	layout := `import React from 'react'

export const metadata = {
  title: 'Next.js App',
  description: 'Boilerplate Next.js App',
}

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
`

	page := `import React from 'react'

export default function Home() {
  return (
    <main style={{ fontFamily: 'system-ui', padding: '4rem' }}>
      <h1>Welcome to Next.js App Router</h1>
      <p>Happy coding!</p>
    </main>
  )
}
`

	appDir := filepath.Join(root, "app")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}

	_ = os.WriteFile(filepath.Join(root, "package.json"), []byte(packageJSON), 0644)
	_ = os.WriteFile(filepath.Join(appDir, "layout.tsx"), []byte(layout), 0644)
	_ = os.WriteFile(filepath.Join(appDir, "page.tsx"), []byte(page), 0644)
	_ = os.WriteFile(filepath.Join(root, "tsconfig.json"), []byte("{}"), 0644)
	return nil
}

func generatePHP(root, name string) error {
	composer := fmt.Sprintf(`{
    "name": "example/%s",
    "description": "Composer PHP Router template",
    "type": "project",
    "require": {
        "php": ">=8.0"
    },
    "autoload": {
        "psr-4": {
            "App\\": "src/"
        }
    }
}
`, name)

	indexPHP := `<?xml version="1.0" encoding="utf-8"?>
<?php
require_once __DIR__ . '/vendor/autoload.php';

$request = $_SERVER['REQUEST_URI'];
$basePath = implode('/', array_slice(explode('/', $_SERVER['SCRIPT_NAME']), 0, -1)) . '/';
$request = str_replace($basePath, '', $request);
$request = '/' . ltrim($request, '/');

switch ($request) {
    case '/':
    case '/index.php':
        header('Content-Type: application/json');
        echo json_encode([
            "message" => "Welcome to PHP Web Router skeleton!",
            "framework" => "Plain PHP Autoloaded"
        ]);
        break;
    default:
        http_response_code(404);
        echo "404 Not Found";
        break;
}
`

	srcDir := filepath.Join(root, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return err
	}

	_ = os.WriteFile(filepath.Join(root, "composer.json"), []byte(composer), 0644)
	_ = os.WriteFile(filepath.Join(root, "index.php"), []byte(indexPHP), 0644)
	_ = os.WriteFile(filepath.Join(srcDir, "User.php"), []byte("<?php\nnamespace App;\nclass User {}\n"), 0644)
	return nil
}

func generateDotnet(root, name string) error {
	programCS := `var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/", () => new { Message = "Hello World from .NET Core Minimal API!", Framework = ".NET 8" });

app.Run();
`

	csproj := `<Project Sdk="Microsoft.NET.Sdk.Web">

  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
    <Nullable>enable</Nullable>
    <ImplicitUsings>enable</ImplicitUsings>
  </PropertyGroup>

</Project>
`

	appsettings := `{
  "Logging": {
    "LogLevel": {
      "Default": "Information",
      "Microsoft.AspNetCore": "Warning"
    }
  },
  "AllowedHosts": "*"
}
`

	_ = os.WriteFile(filepath.Join(root, "Program.cs"), []byte(programCS), 0644)
	_ = os.WriteFile(filepath.Join(root, name+".csproj"), []byte(csproj), 0644)
	_ = os.WriteFile(filepath.Join(root, "appsettings.json"), []byte(appsettings), 0644)
	return nil
}

func generateJava(root, name string) error {
	buildGradle := fmt.Sprintf(`plugins {
    id 'java'
    id 'application'
}

group = 'com.example'
version = '1.0-SNAPSHOT'

repositories {
    mavenCentral()
}

dependencies {
    testImplementation platform('org.junit:junit-bom:5.9.1')
    testImplementation 'org.junit.jupiter:junit-jupiter'
}

application {
    mainClass = 'com.example.%s.Main'
}

test {
    useJUnitPlatform()
}
`, name)

	javaDir := filepath.Join(root, "src", "main", "java", "com", "example", name)
	if err := os.MkdirAll(javaDir, 0755); err != nil {
		return err
	}

	mainJava := fmt.Sprintf(`package com.example.%s;

public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, Developer! This is a plain Java template powered by Gradle build.");
    }
}
`, name)

	_ = os.WriteFile(filepath.Join(root, "build.gradle"), []byte(buildGradle), 0644)
	_ = os.WriteFile(filepath.Join(javaDir, "Main.java"), []byte(mainJava), 0644)
	return nil
}
