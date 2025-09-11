# Latens

> _"Awaken your latent potential"_

**Latens** is a tool that gently awakens sleeping projects, allowing each project to remember and communicate for itself where it left off and what it planned to continue.

## The Problem

As developers, we all have repositories that naturally go into pause:

- Side projects we started with enthusiasm
- POCs and experiments left in light sleep
- Projects that rotated to another team
- Ideas that didn't gain initial traction

When we want to wake them up, we face the fact that they've lost consciousness of their own state:

- The project doesn't remember what it was doing
- It can't explain why it made certain decisions
- It has forgotten at what exact point it fell asleep
- It doesn't know what it planned to achieve

## The Solution

### Latens Score

Visual metric indicating how deep a project's "sleep" is:

- **Light Sleep** (0-30): Recent active development
- **Standard Sleep** (31-60): Moderate pause, intermediate difficult to awake
- **Deep Sleep** (61-100): Deep sleep, requires careful awakening

### Latens Recovery

Process that allows the project to remember its own state:

- Remembers what it was working on
- Identifies technical decisions it had made
- Recognizes the direction it was following
- Recalls its ideas and experiments in progress

### Latens Awakening

The natural awakening of the project, where it speaks for itself:

> _"Mmm... I was working on the authentication API..."_  
> _"I had the JWT integration almost ready..."_  
> _"Next was implementing the middleware..."_  
> _"Oh yes, I needed to test the refresh token..."_

## Tech Stack

- **Frontend**: Next.js 14 + TypeScript + Tailwind CSS
- **Backend**: Go + Fiber + MongoDB
- **Authentication**: GitHub OAuth
- **Analysis**: OpenAI API
- **Deploy**: Vercel + Digital Ocean

## Project Status

🚧 **MVP in development** - Target: September 30, 2025

---

_"Because every project deserves to be awakened, and every developer deserves to remember their dreams."_
