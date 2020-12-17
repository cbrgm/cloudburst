import {css, html, LitElement} from 'lit-element';
import{
    ApiClient,
    TargetsApi,
    InstancesApi
} from '../../api/client/javascript/src/index.js';


class App extends LitElement {
    static get properties() {
        return {
            client: {type: Object},
            scrapeTargets: {type: Array},
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
            
            
                        <div class="card events-card">
                            <header class="card-header">
                                <p class="card-header-title">
                                    BubbleSort-Svc
                                </p>
                            </header>
                            <div class="card-table">
                                <div class="content">
                                    <table class="table is-fullwidth is-striped">
                                        <thead>
                                        <th></th>
                                        <th>Instance</th>
                                        <th>Endpoint</th>
                                        <th>Agent</th>
                                        <th></th>
                                        </thead>
                                        <tbody>
                                        <tr>
                                            <td width="1%"></td>
                                            <td>bubblesort-instance-141231</td>
                                            <td>http://localhost:9090</td>
                                            <td>Fake-Agent</td>
                                            <td class="level-right"><a class="button is-small is-success is-fullwidth" href="#">Running</a></td>
                                        </tr>
                                        <tr>
                                            <td width="1%"></td>
                                            <td>bubblesort-instance-141231</td>
                                            <td>http://localhost:9090</td>
                                            <td>Fake-Agent</td>
                                            <td class="level-right"><a class="button is-small is-info is-fullwidth" href="#">Pending</a></td>
                                        </tr>
                                        <tr>
                                            <td width="1%"></td>
                                            <td>bubblesort-instance-141231</td>
                                            <td>http://localhost:9090</td>
                                            <td>Fake-Agent</td>
                                            <td class="level-right"><a class="button is-small is-info is-loading is-fullwidth" href="#">Progress</a></td>
                                        </tr>
                                        <tr>
                                            <td width="1%"></td>
                                            <td>bubblesort-instance-141231</td>
                                            <td>http://localhost:9090</td>
                                            <td>Fake-Agent</td>
                                            <td class="level-right"><a class="button is-small is-danger is-fullwidth" href="#">Failure</a></td>
                                        </tr>
                                        <tr>
                                            <td width="1%"></td>
                                            <td>bubblesort-instance-141231</td>
                                            <td>http://localhost:9090</td>
                                            <td>Fake-Agent</td>
                                            <td class="level-right"><a class="button is-small is-light is-fullwidth" href="#">Terminated</a></td>
                                        </tr>
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>
            
            
                    </div>
                </div>
            </div>
        `;
    }

    firstUpdated(changedProperties) {

        this.fetchData()

        let instanceEvents = new EventSource(`${this.client.basePath}/instances/events`);
        instanceEvents.onmessage = (event) => console.log(event.data);
    }

    async fetchData() {
        const targets = await new TargetsApi(this.client).listScrapeTargets().then((response) => response);
        const targetNames = targets.map(target => target.name);

        console.log(targetNames)
        targets.forEach(target => {
            new InstancesApi(this.client).getInstances(target.name).then((response) => console.log(response));
        });
    }
}

customElements.define("cloudburst-app", App);