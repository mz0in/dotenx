export async function renderRepoStats(stats) {
    let tableHTML = '';

    for (const [key, value] of Object.entries(stats)) {
        tableHTML += `<tr><td>${key}</td><td>${value}</td></tr>`;
    }

    return `<table>${tableHTML}</table>`;
}
