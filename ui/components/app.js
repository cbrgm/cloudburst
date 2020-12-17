import {css, html, LitElement} from 'lit-element';
import{
    ApiClient,
    TargetsApi,
    InstanceEvent,
} from '../../api/client/javascript/src/index.js';


class App extends LitElement {
    static get properties() {
        return {
            client: {type: Object},
            scrapeTargets: {type: Array},
            eventSource: {type: Object},
        };
    }

    static get styles() {
        return css`
        .columns {
            width: 100%;
            height: 100%;
            margin-left: 0;
        }
        .menu-label {
            color: #8F99A3;
            letter-spacing: 1.3;
            font-weight: 700;
        }
        .menu-list a {
            color: #0F1D38;
            font-size: 14px;
            font-weight: 700;
        }
        .menu-list a:hover {
            background-color: transparent;
            color: #276cda;
        }
        .menu-list a.is-active {
            background-color: transparent;
            color: #276cda;
            font-weight: 700;
        }
        .card {
            box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.18);
            margin-bottom: 2rem;
        }
        .card-header-title {
            color: #8F99A3;
            font-weight: 400;
        }
        .info-tiles {
            margin: 1rem 0;
        }
        .info-tiles .subtitle {
            font-weight: 300;
            color: #8F99A3;
        }
        .hero.welcome.is-info {
            background: #36D1DC;
            background: -webkit-linear-gradient(to right, #5B86E5, #36D1DC);
            background: linear-gradient(to right, #5B86E5, #36D1DC);
        }
        .hero.welcome .title, .hero.welcome .subtitle {
            color: hsl(192, 17%, 99%);
        }
        .card .content {
            font-size: 14px;
        }
        .card-table .table {
            margin-bottom: 0;
        }
        `;
    }

    constructor() {
        super();
        this.client = new ApiClient();
        this.client.basePath = `${window.location.protocol}//${window.location.host}/api/v1`;
        this.scrapeTargets = [];
    }

    render() {
        return html`
            <link rel="stylesheet" href="bulma.min.css">
            <div class="container">
                <div class="columns is-centered">
                    <div class="column is-8">
                        ${this.scrapeTargets.map((scrapeTarget) => html`
                            <div class="card events-card">
                            <header class="card-header">
                                <p class="card-header-title">
                                    ${scrapeTarget.name}
                                </p>
                            </header>
                            
                            <instances-table scrapeTarget="${scrapeTarget.name}"></instances-table>
                            
                            </div>
                        `)}
                    </div>
                </div>
            </div>
        `;
    }

    firstUpdated(changedProperties) {
        new TargetsApi(this.client).listScrapeTargets().then((scrapeTargets) => this.scrapeTargets = scrapeTargets)
    }

}





customElements.define("cloudburst-app", App);