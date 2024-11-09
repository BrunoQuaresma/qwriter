import { input, select, confirm } from "@inquirer/prompts";
import chalk from "chalk";
import { Command } from "commander";
import fs from "node:fs/promises";

const NAME = "qwriter";
const DEFAULT_CONFIG = `./${NAME}.json`;

async function main() {
  const program = new Command();

  program.name(NAME).description("QWriter CLI").version("0.0.1");

  program
    .command("init")
    .description(`Initialize a new ${NAME} config`)
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
