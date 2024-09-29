const { exec } = require('child_process');
const reporter = require('cucumber-html-reporter');
const axios = require('axios');

exec('go test -v', (err, stdout, stderr) => {
    if (err) {
        console.error(`Error executing tests: ${stderr}`);
        return;
    }

    console.log(`Test Results:\n${stdout}`);

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
   

});
