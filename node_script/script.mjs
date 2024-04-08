#!/usr/bin/env node

import fs from "node:fs";
import { writeFile } from "node:fs/promises";
import readline from "node:readline/promises";
import { PerformanceObserver, performance } from "node:perf_hooks";

const perf = new PerformanceObserver((list) => {
  for (const entry of list.getEntries()) {
    console.log(entry.duration / 1000, "sec");
  }
});

perf.observe({ entryTypes: ["measure"], buffered: true });

(async function main() {
  try {
    performance.mark("start");

    let data = [];

    const rl = readline.createInterface({
      input: fs.createReadStream("./us_cities_states_counties.csv"),
    });

    const clean = async (msg) => {
      if (msg === "close") {
        const cities = data
          .map((d) => {
            const row = d.split(",");
            return {
              city: row[0],
              state_short: row[1],
              state_full: row[2],
            };
          })
          .reduce((acc, cur) => {
            if (!acc.some((d) => d.city.includes(cur.city))) {
              return [...acc, cur];
            }
            return acc;
          }, [])
          .sort((a, b) => a.city.localeCompare(b.city));

        const states = data
          .reduce((acc, cur) => {
            const row = cur.split(",").slice(1, 3).join(",");
            if (Array.isArray(acc) && !acc.includes(row)) {
              return [...acc, row];
            }
            return [...acc];
          }, [])
          .map((d) => {
            const row = d.split(",");
            return {
              state_short: row[0],
              state_full: row[1],
            };
          })
          .sort((a, b) => a.state_short.localeCompare(b.state_short));

        Promise.all([
          await writeFile("cities.json", JSON.stringify(cities)),
          await writeFile("states.json", JSON.stringify(states)),
        ]);

        console.log(cities.length);
        console.log(states.length);

        performance.mark("end");
        performance.measure("processing", "start", "end");
      }
    };

    rl.on("line", (ln) => {
      const d = ln.toString().trim();
      if (d === "city|state_short|state_full|county|city_alias") return;
      const city = d.split("|").slice(0, -2).join(",");
      if (data.includes(city)) return;
      data.push(city);
    });

    rl.on("close", () => clean("close"));
  } catch (e) {
    console.log(e);
  }
})();
