async function loadCases() {

    const res = await fetch("/api/cases");
    const cases = await res.json();

    const container = document.getElementById("casesContainer");

    container.innerHTML = "";

    cases.forEach(c => {

        const percent = c.target_amount > 0
            ? (c.amount_raised / c.target_amount) * 100
            : 0;

        container.innerHTML += `
        <article class="case-card">

            <div class="card-details-panel">

                <h2>${c.title}</h2>

                <p>${c.description}</p>

                <p>Target: Ksh ${c.target_amount}</p>
                <p>Raised: Ksh ${c.amount_raised}</p>

                <div style="background:#eee;height:10px;border-radius:5px;">
                    <div style="width:${percent}%;background:green;height:10px;border-radius:5px;"></div>
                </div>

                <p>${percent.toFixed(1)}%</p>

                <div class="card-actions">
                    <a href="case-details.html?id=${c.id}">View</a>
                    <a href="donate.html?id=${c.id}">Donate</a>
                </div>

            </div>

        </article>
        `;
    });
}

loadCases();