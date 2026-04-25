---
description: Show how to create a new Claude skill/command for this project
---

## How to create a skill

Skills are markdown files in `.claude/commands/`. They become `/command-name` slash commands.

### File format

```markdown
---
description: One-line description shown in the skill list
---

Your prompt here. Claude reads this and executes it.

Use $ARGUMENTS to capture what the user types after the command.
Example: /my-skill some text  →  $ARGUMENTS = "some text"
```

### Create a new skill

To create a skill called `/my-skill`:

1. Create `.claude/commands/my-skill.md`
2. Add frontmatter + prompt
3. Use it immediately — no restart needed

### Examples in this project

| Command | File | Purpose |
|---------|------|---------|
| `/backend` | commands/backend.md | Go backend tasks |
| `/frontend` | commands/frontend.md | React tasks |
| `/k8s` | commands/k8s.md | Kubernetes manifests |
| `/jenkins` | commands/jenkins.md | CI/CD pipelines |
| `/docs` | commands/docs.md | Documentation |
| `/status` | commands/status.md | Check pod health |
| `/logs` | commands/logs.md | Tail service logs |
| `/deploy` | commands/deploy.md | Redeploy a service |

### Tips

- Keep the prompt focused — one skill, one job
- Use `$ARGUMENTS` when the skill needs user input
- Reference project context (ports, namespaces, paths) directly in the prompt
- Skills are committed to git — the whole team shares them

Now tell me what skill you want to create and I'll build it.
$ARGUMENTS
