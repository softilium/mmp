import moment from 'moment/min/moment-with-locales';

export const glob = {

  fmtDate: (date) => {
    moment.locale('ru');
    return moment(date).fromNow();
  }

};
