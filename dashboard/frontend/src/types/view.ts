interface Param {
  type: string;
  operator: string;
  key: string;
  value: ParamValue;
}

interface ParamValue {
  item: string;
  list: string[];
}

interface Form {
  [key: string]: Date | string | string[] | null;
}

interface SearchParam {
  key: string;
  name: string;
  type: string;
}

export {
  Param, Form, ParamValue, SearchParam,
};
