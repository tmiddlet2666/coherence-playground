package com.tmiddlet.demo;

import java.util.List;

import org.springframework.ai.support.ToolCallbacks;
import org.springframework.ai.tool.ToolCallback;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean; // Needed to expose the tools

@SpringBootApplication
public class MCPCoherenceApplication {

	public static void main(String[] args) {
		SpringApplication.run(MCPCoherenceApplication.class, args);
	}

	@Bean
	public List<ToolCallback> tools(MCPCoherence mcpCoherence) {
		return List.of(ToolCallbacks.from(mcpCoherence));
	}
}
