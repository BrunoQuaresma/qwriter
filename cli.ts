import { input, select, confirm } from "@inquirer/prompts";
import chalk from "chalk";
import { Command } from "commander";
import fs from "node:fs/promises";
import { glob } from "glob";
import OpenAI from "openai";
import ora from "ora";

const DEFAULT_CONFIG = `./qwriter.json`;

async function main() {
  const program = new Command();

  program
    .name("QWriter")
    .description(
      "QWriter is a CLI tool designed to enhance and polish the copy in your text and code files. It reviews content for tone, clarity, and readability, making it easy to improve everything from documentation and comments to marketing copy. With configurable options, QWriter tailors its suggestions to your specific audience and goals, giving your content a professional edge."
    )
    .version("0.0.1");

  program
    .command("init")
    .description(
      "Initializes a QWriter configuration file with your preferred settings for generating and formatting content. This command will guide you through setting up the basic options to get started quickly."
    )
    .action(async () => {
      const goal = await input({
        message: chalk.bold("What is your main objective?"),
      });
      const topic = await input({
        message: chalk.bold("What is the main topic or focus of your content?"),
      });
      const audience = await input({
        message: chalk.bold("Who is your target audience for this content?"),
      });
      const tone = await select({
        message: chalk.bold("What tone would you like for the content?"),
        choices: [
          {
            name: "Formal",
            value: "formal",
            description:
              "Professional and respectful, suitable for business or academic contexts.",
          },
          {
            name: "Casual",
            value: "casual",
            description:
              "Conversational and relaxed, ideal for engaging a general audience.",
          },
          {
            name: "Persuasive",
            value: "persuasive",
            description:
              "Compelling and motivating, perfect for marketing or calls-to-action.",
          },
          {
            name: "Instructional",
            value: "instructional",
            description:
              "Clear and directive, suitable for step-by-step guides or tutorials.",
          },
          {
            name: "Friendly",
            value: "friendly",
            description:
              "Warm and approachable, helping to create a relatable and welcoming tone.",
          },
          {
            name: "Technical",
            value: "technical",
            description:
              "Precise and detail-oriented, ideal for expert audiences and documentation.",
          },
        ],
      });
      const configPath = await input({
        message: chalk.bold(
          "Where would you like to save the configuration file?"
        ),
        default: DEFAULT_CONFIG,
      });

      // Check if the config file already exists
      let [_, err] = await catchError(fs.access(configPath, fs.constants.F_OK));
      if (!err) {
        const overwrite = await confirm({
          message: `The configuration file already exists at ${configPath}. Do you want to overwrite it?`,
        });
        if (!overwrite) {
          console.log(`${chalk.red("✖")} Configuration file not saved.`);
          return;
        }
      }

      // Save the configuration to the file
      [_, err] = await catchError(
        fs.writeFile(
          configPath,
          JSON.stringify({ goal, topic, audience, tone }, null, 2)
        )
      );
      if (err) {
        console.log(`${chalk.red("✖")} Failed to save configuration file.`);
        console.error(chalk.red(err));
        return;
      }
      console.log(
        `${chalk.green(
          "✔"
        )} Configuration file saved successfully to ${chalk.bold(configPath)}`
      );
    });

  program
    .command("review")
    .description(
      "Reviews and rewrites content in the specified file path, enhancing tone, clarity, and structure based on your QWriter settings. Use this command to refine multiple files, making your content more effective and engaging."
    )
    .argument("<path>", "The file path to the content to review.")
    .option(
      "-c, --config <path>",
      "The path to the QWriter configuration file.",
      DEFAULT_CONFIG
    )
    .option(
      "-x, --context <context>",
      "Allows you to provide additional context to guide QWriter's content improvements. Use this to specify extra details about your project, audience, or unique requirements that can help tailor the output more precisely.",
      ""
    )
    .action(async (path, options: { config: string; context: string }) => {
      let [_, err] = await catchError(
        fs.access(options.config, fs.constants.F_OK)
      );
      if (err) {
        console.log(
          `${chalk.red("✖")} No configuration exists at ${chalk.bold(
            options.config
          )}.`
        );
        return;
      }

      const [config, readErr] = await catchError(
        fs.readFile(options.config, "utf-8")
      );
      if (readErr) {
        console.log(`${chalk.red("✖")} Failed to read configuration file.`);
        console.error(chalk.red(readErr));
        return;
      }
      if (!config) {
        console.log(`${chalk.red("✖")} Empty configuration file.`);
        return;
      }
      const { goal, topic, audience, tone } = JSON.parse(config) as {
        goal: string;
        topic: string;
        audience: string;
        tone: string;
      };

      const prompt = `
        You are an expert editor skilled in enhancing content for clarity, engagement, and effectiveness across various formats.
        
        Objective: The user's goal is to "${goal}".
        
        Content Description:
        - Topic: The content is about "${topic}".
        - Audience: The intended audience is "${audience}".
        - Tone: The content should use a "${tone}" tone to resonate with the target readers.
        
        Additional Context:
        ${options.context || "No additional context provided."}
        
        Instructions:
        1. For content targeted at a general audience, ensure language is clear, friendly, and avoids unnecessary jargon.
        2. If the content is technical (e.g., documentation or code comments), focus on simplifying complex explanations while maintaining accuracy. Avoid altering code functionality if present.
        3. Adapt the language to fit the "${tone}" tone specified, making it ${tone.toLowerCase()}.
        
        Enhance the text as follows:
        - Adjust the language to be suitable for the "${audience}" audience.
        - Improve readability, conciseness, and flow, aligning with the specified "${goal}".
        - Where necessary, rephrase for engagement and clarity.
        - Ensure terminology aligns with the topic while being accessible to the intended audience.
        
        Output only the revised text without extra explanations or formatting, so it can be seamlessly integrated into the users original files.
      `;

      if (!process.env.OPENAI_API_KEY) {
        console.log(
          `${chalk.red("✖")} Missing OPENAI_API_KEY environment variable.`
        );
        return;
      }
      const openAI = new OpenAI();

      const files = await glob(path);
      for (const file of files) {
        const [content, readErr] = await catchError(fs.readFile(file, "utf-8"));
        if (readErr) {
          console.log(
            `${chalk.red("✖")} Failed to read content from ${chalk.bold(file)}.`
          );
          console.error(chalk.red(readErr));
          continue;
        }
        if (!content) {
          console.log(
            `${chalk.red("✖")} Empty content in ${chalk.bold(file)}.`
          );
          continue;
        }

        const spinner = ora(`Reviewing ${chalk.bold(file)}...`).start();
        const [chatResult, chatError] = await catchError(
          openAI.chat.completions.create({
            messages: [
              { role: "system", content: prompt },
              { role: "user", content },
            ],
            model: "gpt-4o",
          })
        );
        spinner.stop();
        if (chatError) {
          console.log(
            `${chalk.red("✖")} Failed to review content in ${chalk.bold(file)}.`
          );
          console.error(chalk.red(chatError));
          continue;
        }
        const revisedContent = chatResult!.choices[0].message.content;
        if (!revisedContent) {
          console.log(
            `${chalk.red("✖")} No revised content generated for ${chalk.bold(
              file
            )}.`
          );
          continue;
        }
        const [_, writeErr] = await catchError(
          fs.writeFile(file, revisedContent)
        );
        if (writeErr) {
          console.log(
            `${chalk.red("✖")} Failed to write revised content to ${chalk.bold(
              file
            )}.`
          );
          console.error(chalk.red(writeErr));
          continue;
        }
        console.log(
          `${chalk.green("✔")} Revised content saved to ${chalk.bold(file)}`
        );
      }
    });

  program.parse();
}

main();

// Utility functions

type Success<T> = [T, undefined];
type Failure = [undefined, unknown];

async function catchError<T>(
  promise: Promise<T>
): Promise<Success<T> | Failure> {
  try {
    const result = await promise;
    return [result, undefined];
  } catch (error) {
    return [undefined, error];
  }
}
