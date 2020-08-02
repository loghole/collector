function FilterTags(data: string[], selected: string[], input = ''): string[] {
  const tmp = data.filter(
    (option: string) => option
      .toString()
      .toLowerCase()
      .indexOf(input.toLowerCase()) >= 0,
  );

  return tmp.filter((option: string) => !selected.includes(option));
}

export default FilterTags;
