// tslint:disable-next-line:no-var-requires
const spinner = require<string>("assets/spinner.gif");
// tslint:disable-next-line:no-var-requires
const noResults = require<string>("assets/noresults.svg");

import _ from "lodash";
// import moment from "moment";
// import { Line } from "rc-progress";
import React from "react";
import { connect } from "react-redux";

import * as protos from "src/js/protos";
import { queriesKey, refreshQueries } from "src/redux/apiReducers";
import { LocalSetting } from "src/redux/localsettings";
import { AdminUIState } from "src/redux/state";
import { TimestampToMoment } from "src/util/convert";
import docsURL from "src/util/docs";
import Dropdown, { DropdownOption } from "src/views/shared/components/dropdown";
import Loading from "src/views/shared/components/loading";
import { PageConfig, PageConfigItem } from "src/views/shared/components/pageconfig";
import { SortSetting } from "src/views/shared/components/sortabletable";
import { ColumnDescriptor, SortedTable } from "src/views/shared/components/sortedtable";
import { ToolTipWrapper } from "src/views/shared/components/toolTip";

type Query = protos.cockroach.server.serverpb.QueriesResponse.Query;

type JobType = protos.cockroach.sql.jobs.Type;
const jobType = protos.cockroach.sql.jobs.Type;

const QueriesRequest = protos.cockroach.server.serverpb.QueriesRequest;

const statusOptions = [
  { value: "", label: "All" },
  { value: "pending", label: "Pending" },
  { value: "running", label: "Running" },
  { value: "paused", label: "Paused" },
  { value: "canceled", label: "Canceled" },
  { value: "succeeded", label: "Succeeded" },
  { value: "failed", label: "Failed" },
];

const statusSetting = new LocalSetting<AdminUIState, string>(
  "jobs/status_setting", s => s.localSettings, statusOptions[0].value,
);

const typeOptions = [
  { value: jobType.UNSPECIFIED.toString(), label: "All" },
  { value: jobType.BACKUP.toString(), label: "Backups" },
  { value: jobType.RESTORE.toString(), label: "Restores" },
  { value: jobType.SCHEMA_CHANGE.toString(), label: "Schema Changes" },
];

const typeSetting = new LocalSetting<AdminUIState, number>(
  "jobs/type_setting", s => s.localSettings, jobType.UNSPECIFIED,
);

const showOptions = [
  { value: "50", label: "Latest 50" },
  { value: "0", label: "All" },
];

const showSetting = new LocalSetting<AdminUIState, string>(
  "jobs/show_setting", s => s.localSettings, showOptions[0].value,
);

// Moment cannot render durations (moment/moment#1048). Hack it ourselves.
// const formatDuration = (d: moment.Duration) =>
//   [Math.floor(d.asHours()).toFixed(0), d.minutes(), d.seconds()]
//     .map(c => ("0" + c).slice(-2))
//     .join(":");

// class QueryStatusCell extends React.Component<{ query: Query }, {}> {
//   is(...statuses: string[]) {
//     // return statuses.indexOf(this.props.query.status) !== -1;
//     return statuses.indexOf("Pending") !== -1;
//   }

//   renderProgress() {
//     if (this.is("succeeded", "failed", "canceled")) {
//       // return <div className="jobs-table__status">{this.props.query.status}</div>;
//       return <div className="jobs-table__status">succeeded</div>;
//     }
//     // const percent = this.props.query.fraction_completed * 100;
//     const percent = 50;
//     return <div>
//       <Line percent={percent} strokeWidth={10} trailWidth={10} className="jobs-table__progress-bar" />
//       <span title={percent.toFixed(3) + "%"}>{percent.toFixed(1) + "%"}</span>
//     </div>;
//   }

//   renderDuration() {
//     const started = TimestampToMoment(this.props.query.start);
//     const finished = TimestampToMoment(this.props.query.start);
//     const modified = TimestampToMoment(this.props.query.start);
//     if (this.is("pending", "paused")) {
//       return _.capitalize("<STATUS>"); // this.props.query.status
//     } else if (this.is("running")) {
//       const fractionCompleted = 0.5;
//       if (fractionCompleted > 0) {
//         const duration = modified.diff(started);
//         const remaining = duration / fractionCompleted - duration;
//         return formatDuration(moment.duration(remaining)) + " remaining";
//       }
//     } else if (this.is("succeeded")) {
//       return "Duration: " + formatDuration(moment.duration(finished.diff(started)));
//     }
//   }

//   render() {
//     return <div>
//       {this.renderProgress()}
//       <span className="jobs-table__duration">{this.renderDuration()}</span>
//     </div>;
//   }
// }

// Specialization of generic SortedTable component:
//   https://github.com/Microsoft/TypeScript/issues/3960
//
// The variable name must start with a capital letter or JSX will not recognize
// it as a component.
// tslint:disable-next-line:variable-name
const QueriesSortedTable = SortedTable as new () => SortedTable<Query>;

const QueriesTableColumns: ColumnDescriptor<Query>[] = [
  {
    title: "Query ID",
    cell: query => query.query_id,
    sort: query => query.query_id,
  },
  {
    title: "NodeID",
    cell: query => String(query.node_id),
    sort: query => query.node_id,
  },
  {
    title: "Query",
    cell: query => query.query,
    sort: query => query.query,
  },
  {
    title: "Client Address",
    cell: query => query.client_address,
    sort: query => query.client_address,
  },
  {
    title: "Application Name",
    cell: query => query.application_name,
    sort: query => query.application_name,
  },
  {
    title: "User",
    cell: query => query.username,
    sort: query => query.username,
  },
  {
    title: "Start Time",
    cell: query => TimestampToMoment(query.start).fromNow(),
    sort: query => TimestampToMoment(query.start).valueOf(),
  },
  // {
  //   title: "Status",
  //   cell: query => <QueryStatusCell query={query} />,
  //   sort: query => query.fraction_completed,
  // },
];

const sortSetting = new LocalSetting<AdminUIState, SortSetting>(
  "jobs/sort_setting",
  s => s.localSettings,
  { sortKey: 2 /* creation time */, ascending: false },
);

interface QueriesTableProps {
  sort: SortSetting;
  status: string;
  show: string;
  type: number;
  setSort: (value: SortSetting) => void;
  setStatus: (value: string) => void;
  setShow: (value: string) => void;
  setType: (value: JobType) => void;
  refreshQueries: typeof refreshQueries;
  queries: Query[];
  queriesValid: boolean;
}

const titleTooltip = (
  <span>
    Some jobs can be paused or canceled through SQL. For details, view the docs
    on the <a href={docsURL("pause-job.html")}><code>PAUSE JOB</code></a> and <a
      href={docsURL("cancel-job.html")}><code>CANCEL JOB</code></a> statements.
  </span>
);

class QueriesTable extends React.Component<QueriesTableProps, {}> {
  static title() {
    return (
      <div>
        Queries
        <div className="section-heading__tooltip">
          <ToolTipWrapper text={titleTooltip}>
            <div className="section-heading__tooltip-hover-area">
              <div className="section-heading__info-icon">!</div>
            </div>
          </ToolTipWrapper>
        </div>
      </div>
    );
  }

  refresh(props = this.props) {
    props.refreshQueries(new QueriesRequest({
      status: props.status,
      type: props.type,
      limit: parseInt(props.show, 10),
    }));
  }

  componentWillMount() {
    this.refresh();
  }

  componentWillReceiveProps(props: QueriesTableProps) {
    this.refresh(props);
  }

  onStatusSelected = (selected: DropdownOption) => {
    this.props.setStatus(selected.value);
  }

  onTypeSelected = (selected: DropdownOption) => {
    this.props.setType(parseInt(selected.value, 10));
  }

  onShowSelected = (selected: DropdownOption) => {
    this.props.setShow(selected.value);
  }

  render() {
    const data = this.props.queries && this.props.queries.length > 0 && this.props.queries;
    return <div className="jobs-page">
      <div>
        <PageConfig>
          <PageConfigItem>
            <Dropdown
              title="Status"
              options={statusOptions}
              selected={this.props.status}
              onChange={this.onStatusSelected}
            />
          </PageConfigItem>
          <PageConfigItem>
            <Dropdown
              title="Type"
              options={typeOptions}
              selected={this.props.type.toString()}
              onChange={this.onTypeSelected}
            />
          </PageConfigItem>
          <PageConfigItem>
            <Dropdown
              title="Show"
              options={showOptions}
              selected={this.props.show}
              onChange={this.onShowSelected}
            />
          </PageConfigItem>
        </PageConfig>
      </div>
      <Loading loading={_.isNil(this.props.queries)} className="loading-image loading-image__spinner" image={spinner}>
        <Loading loading={_.isEmpty(data)} className="loading-image" image={noResults}>
          <section className="section">
            <QueriesSortedTable
              data={data}
              sortSetting={this.props.sort}
              onChangeSortSetting={this.props.setSort}
              className="jobs-table"
              columns={QueriesTableColumns}
            />
          </section>
        </Loading>
      </Loading>
    </div>;
  }
}

const mapStateToProps = (state: AdminUIState) => {
  const sort = sortSetting.selector(state);
  const status = statusSetting.selector(state);
  const show = showSetting.selector(state);
  const type = typeSetting.selector(state);
  const key = queriesKey(status, type, parseInt(show, 10));
  const queries = state.cachedData.queries[key];
  return {
    sort, status, show, type,
    queries: queries && queries.data && queries.data.queries,
    queriesValid: queries && queries.valid,
  };
};

const actions = {
  setSort: sortSetting.set,
  setStatus: statusSetting.set,
  setShow: showSetting.set,
  setType: typeSetting.set,
  refreshQueries,
};

export default connect(mapStateToProps, actions)(QueriesTable);
