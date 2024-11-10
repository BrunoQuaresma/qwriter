# QWriter

QWriter is a CLI tool that goes beyond “just another ChatGPT API wrapper” and makes your text and code files shine. Sure, it leverages the ChatGPT API, but this isn’t about simple questions and answers. QWriter is your trusty editing sidekick, making technical documentation readable, code comments crystal-clear, and marketing copy engaging.

**Why QWriter?** Because even the best writers need an editor, and let’s face it—sometimes your comments need more than a sprinkle of polish. With QWriter, every file can get that professional, AI-powered touch.

## Quick Start

1. Set up your OpenAI API key as an environment variable:

   ```bash
   export OPENAI_API_KEY=your_openai_api_key_here
   ```

2. Run QWriter with `npx` (no installation required):

   ```bash
   npx qwriter init
   ```

   This will launch the `init` command to set up your configuration file.

Now you’re ready to use QWriter commands directly with `npx qwriter <command>`!

## Commands & Examples

### 1. `qwriter init` - Initialize your config file

First, set up a config file to let QWriter know your objective, topic, audience, and preferred tone. QWriter will prompt you with easy questions to help you configure it.

#### Example

```bash
npx qwriter init
```

**Sample Output:**

```plaintext
✔ What is your main objective? Improve clarity and tone of product copy.
✔ What is the main topic or focus of your content? A new feature for SaaS billing.
✔ Who is your target audience for this content? Small business owners.
✔ What tone would you like for the content? Friendly
✔ Where would you like to save the configuration file? ./qwriter.json
```

> _"Configuration file saved successfully to ./qwriter.json"_ 🎉

---

### 2. `qwriter review` - Review and Rewrite Content

With your config file set up, use this command to review and rewrite files, whether it’s for clear documentation, engaging marketing copy, or understandable code comments.

#### Example Usage Scenarios

---

#### Scenario 1: Sprucing Up Marketing Copy

_Your SaaS feature release sounds stiff. QWriter can help make it more friendly!_

```bash
npx qwriter review ./content/marketing-copy.txt -c ./qwriter.json
```

---

#### Scenario 2: Making Technical Documentation Accessible

_Turn a dense API setup guide into something beginners won’t fear!_

```bash
npx qwriter review ./docs/api-setup.md -c ./qwriter.json -x "Make it approachable for beginners"
```

---

#### Scenario 3: Cleaning Up Code Comments for Clarity

_Your code comments need translation. Let QWriter help._

```bash
npx qwriter review ./src/new-feature.js -c ./qwriter.json -x "Clarify comments without altering functionality"
```

---

## File Content Examples (Before & After)

Want to see the magic in action? Here’s a peek at what QWriter can do:

---

### Example 1: Marketing Copy (`marketing-copy.txt`)

**Before:**

```plaintext
Our new automated billing feature is designed to streamline your payment processes. It ensures accurate, timely invoicing and provides detailed reports. Manage your finances with ease and gain insights into your business performance.
```

**After QWriter:**

```plaintext
Meet our new automated billing feature—built to simplify your payments! Now, you can automate invoicing, receive timely reports, and get clear insights into your business performance, so you can focus on what matters most.
```

---

### Example 2: Technical Documentation (`api-setup.md`)

**Before:**

```markdown
# API Setup Guide

1. **Obtain API Key**: Go to the Developer Portal and generate an API key. This key will allow you to authenticate requests.
2. **Install SDK**: Use the command `npm install our-sdk` to install the SDK for JavaScript.
3. **Configure SDK**: Initialize the SDK by passing your API key into the configuration function.
4. **Make Your First Request**: Use the SDK’s `fetchData` method to retrieve your data. Refer to the documentation for endpoint details.
```

**After QWriter:**

```markdown
# API Setup Guide

1. **Get Your API Key**: Visit our Developer Portal to generate an API key, which will authenticate your requests.
2. **Install the SDK**: Run `npm install our-sdk` to add the SDK to your project.
3. **Set Up the SDK**: Pass your API key to initialize the SDK.
4. **Retrieve Data**: Call the SDK’s `fetchData` method to pull data from our API. See the documentation for details.
```

---

### Example 3: Code Comments (`new-feature.js`)

**Before:**

```javascript
// Function to process user data and calculate metrics
function processData(data) {
  // This part loops through the data to calculate average metrics
  for (let i = 0; i < data.length; i++) {
    // Calculate average and add it to the result array
    const average = (data[i].value1 + data[i].value2) / 2;
    result.push(average);
  }
}

// A function to log the results to the console
// This helps us verify the calculations
function logResults(results) {
  console.log("Results:", results);
}
```

**After QWriter:**

```javascript
// Processes user data to calculate average metrics
function processData(data) {
  // Loop through each data point and calculate the average of values
  for (let i = 0; i < data.length; i++) {
    const average = (data[i].value1 + data[i].value2) / 2;
    result.push(average); // Add the average to results
  }
}

// Logs the calculated results to the console for verification
function logResults(results) {
  console.log("Results:", results);
}
```

---

## Command Options

- **`-c, --config <path>`**: Specify a custom config file path (defaults to `./qwriter.json`).
- **`-x, --context <context>`**: Provide additional context to guide QWriter’s content improvements (e.g., "tailor to non-technical readers").

## License

MIT License

## Contributing

We’re open to PRs and issues! QWriter is always in beta in our hearts, and we love fresh ideas.

---

And that’s QWriter—your friendly, AI-powered editor for all things text and code. Because even your files deserve a little polish! 🧑‍💻✨
