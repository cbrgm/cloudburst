import {css, html, LitElement} from 'lit-element';
import{
    ApiClient,
    InstancesApi,
    InstanceEvent,
    Instance,
} from '../../api/client/javascript/src/index.js';


class InstancesTable extends LitElement {
    static get properties() {
        return {
            client: {type: Object},

            scrapeTarget: {type: String},
            instances: {type: Array},
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

        this.scrapeTarget = "";
        this.instances = [];
    }

    render() {
        return html`
            <link rel="stylesheet" href="bulma.min.css">
            
            ${this.instances.length === 0 ? '' : html`
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
                        ${this.instances.map((instance) => html`
                        <tr>
                           <td width="1%"></td>
                           <td>${instance.name}</td>
                           <td>${instance.endpoint}</td>
                           <td>${instance.status.agent}</td>
                           ${instance.status.status === "pending" ? html`<td class="level-right"><a class="button is-small is-rounded is-fullwidth" href="#">Pending</a></td>` : ''}
                           ${instance.status.status === 'progress' ? html`<td class="level-right"><a class="button is-small is-rounded is-info is-loading is-fullwidth" href="#">Progress</a></td>` : ''}
                           ${instance.status.status === 'running' ? html`<td class="level-right"><a class="button is-small is-rounded is-success is-fullwidth" href="#">Running</a></td>` : ''}
                           ${instance.status.status === 'failure' ? html`<td class="level-right"><a class="button is-small is-rounded is-danger is-fullwidth" href="#">Failure</a></td>` : ''}
                           ${instance.status.status === 'terminated' ? html`<td class="level-right"><a class="button is-small is-rounded is-light is-fullwidth" href="#">Terminated</a></td>` : ''}
                        </tr>
                        `)}
                     </tbody>
                  </table>
               </div>
            </div>
            `}
        `;
    }

    firstUpdated(changedProperties) {

        let instanceEvents = new EventSource(`${this.client.basePath}/instances/events`);

        new InstancesApi(this.client).getInstances(this.scrapeTarget).then((instances) => this.instances = instances);
        instanceEvents.onmessage = (event) => this.updateInstances(event);
    }

    updateInstances(event) {
        let instanceEvent = InstanceEvent.constructFromObject(JSON.parse(event.data));

        if (this.scrapeTarget != instanceEvent.target) {
            return
        }

        if (instanceEvent.type === "save") {
            let next = [...this.instances]
            instanceEvent.data.forEach((item) => pushToArray(next, item))
            this.instances = next
        }

        if (instanceEvent.type === "remove") {
            let next = [...this.instances]
            for (let i = 0; i < instanceEvent.data.length; i++) {
                let toRemove = instanceEvent.data[i]
                let index = next.findIndex((e) => e.name === toRemove.name);
                if (index > -1) {
                    next.splice(index, 1)
                }
            }
            this.instances = next
        }
    }
}

function pushToArray(arr, obj) {
    const index = arr.findIndex((e) => e.name === obj.name);
    if (index === -1) {
        arr.push(obj);
    } else {
        arr[index] = obj;
    }
}


customElements.define("instances-table", InstancesTable);