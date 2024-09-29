const { exec } = require('child_process');
const reporter = require('cucumber-html-reporter');
const axios = require('axios');

// Define your Slack webhook URL
const webhookUrl = 'https://hooks.slack.com/services/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'; // Replace with your actual webhook URL

// Function to send results to Slack
function sendResultsToSlack(results) {
    const reportLink = 'http://yourserver.com/path/to/test/report/cucumber_report.html'; // Replace with the actual URL to your HTML report
    const message = {
        text: `*Test Results:* \n*Passed Scenarios:* \n\`\`\`${results}\`\`\`\n\n*For detailed results, check the [HTML report](${reportLink}).`
    };

    axios.post(webhookUrl, message)
        .then(response => {
            console.log('Message posted to Slack');
        })
        .catch(error => {
            console.error(`Error posting to Slack: ${error.response ? error.response.data : error.message}`);
        });
}

// Step 1: Run Go tests
exec('go test -v', (err, stdout, stderr) => {
    if (err) {
        console.error(`Error executing tests: ${stderr}`);
        return;
    }

    console.log(`Test Results:\n${stdout}`);



    // Step 2: Generate HTML report
    const options = {
        theme: 'bootstrap',
        jsonFile: 'test/report/cucumber_report.json', // Ensure this matches your output path
        output: 'test/report/cucumber_report.html', // Path where the HTML report will be saved
        reportSuiteAsScenarios: true,
        launchReport: true,
        metadata: {
            "App Version": "0.0.0",
            "Test Environment": "STAGING",
            "Browser": "Chrome",
            "Platform": "Windows",
            "Parallel": "Scenarios",
            "Executed": "Remote"
        }
    };

    reporter.generate(options);

    console.log('HTML report generated: test/report/cucumber_report.html');

    // Step 3: Send results to Slack
    sendResultsToSlack(stdout);
});
