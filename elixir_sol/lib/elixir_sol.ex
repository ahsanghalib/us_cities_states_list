NimbleCSV.define(CSV, separator: "|", escape: "\"")

defmodule Benchmark do
  def measure(function) do
    time =
      function
      |> :timer.tc()
      |> elem(0)
      |> Kernel./(1_000_000)
      |> to_string()

    IO.puts(time <> " sec")
  end
end

defmodule Solution do
  def process do
    list =
      "us_cities_states_counties.csv"
      |> File.stream!()
      |> CSV.parse_stream()

    cities =
      list
      |> Stream.map(fn [city, state_short, state_full, _, _] ->
        [city, state_short, state_full]
      end)
      |> Enum.uniq()
      |> Enum.sort()
      |> Enum.map(fn [city, state_short, state_full] ->
        %{
          "city" => city,
          "state_short" => state_short,
          "state_full" => state_full
        }
      end)

    states =
      list
      |> Stream.map(fn [_, state_short, state_full, _, _] ->
        [state_short, state_full]
      end)
      |> Enum.uniq()
      |> Enum.sort()
      |> Enum.map(fn [state_short, state_full] ->
        %{
          "state_short" => state_short,
          "state_full" => state_full
        }
      end)

    IO.puts(length(cities))
    IO.puts(length(states))

    states_file = fn d -> File.write!("states.json", d) end
    cities_file = fn d -> File.write!("cities.json", d) end

    states |> Poison.encode!() |> states_file.()
    cities |> Poison.encode!() |> cities_file.()
  end

  def run do
    Benchmark.measure(&Solution.process/0)
  end
end
