#!/usr/bin/env elixir

Mix.install([:poison])

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
  def process() do
    list =
      File.stream!("./us_cities_states_counties.csv")
      |> Stream.map(&String.trim(&1))
      |> Stream.map(&String.split(&1, "|"))
      |> Stream.filter(fn
        ["city" | _] -> false
        _ -> true
      end)

    cities =
      list
      |> Stream.map(&(Stream.drop(&1, -2) |> Enum.to_list()))
      |> Enum.uniq()
      |> Enum.map(fn e ->
        %{
          city: Enum.at(e, 0),
          state_short: Enum.at(e, 1),
          state_full: Enum.at(e, 2)
        }
      end)

    states =
      list
      |> Stream.map(fn l -> l |> Enum.to_list() |> Enum.slice(1, 2) end)
      |> Enum.uniq()
      |> Enum.map(fn e ->
        %{
          state_short: Enum.at(e, 0),
          state_full: Enum.at(e, 1)
        }
      end)

    IO.inspect(length(cities))
    IO.inspect(length(states))

    states_file = fn d -> File.write!("states.json", d) end
    cities_file = fn d -> File.write!("cities.json", d) end

    states |> Poison.encode!() |> states_file.()
    cities |> Poison.encode!() |> cities_file.()
  end
end

Benchmark.measure(&Solution.process/0)
