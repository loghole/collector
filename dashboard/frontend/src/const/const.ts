const OperatorMultiple = ['=', '!=', 'LIKE', 'NOT LIKE'] as string[];

const SingleParam = ['<', '<=', '>', '>='] as string[];

const OperatorsSingle = OperatorMultiple.concat(SingleParam) as string[];

const IntervalRegexp = new RegExp('^([0-9]+)(s|sec|m|min|h|hr|hour|d|day)?$', 'i');

export {
  OperatorsSingle,
  OperatorMultiple,
  SingleParam,
  IntervalRegexp,
};
