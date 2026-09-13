#!/usr/bin/env node

/**
 * Script to add // @ts-nocheck directive to Orval-generated files
 * This script is called by Orval's afterAllFilesWrite hook
 */

import fs from 'fs'
import path from 'path'

const TS_NOCHECK_DIRECTIVE = '// @ts-nocheck'

function addTsNocheckToFile(filePath) {
	try {
		const content = fs.readFileSync(filePath, 'utf8')

		// Check if file already has @ts-nocheck directive
		if (content.startsWith(TS_NOCHECK_DIRECTIVE)) {
			return false // Already has directive
		}

		// Check if file is a TypeScript file
		if (!filePath.endsWith('.ts')) {
			return false // Not a TypeScript file
		}

		// Add @ts-nocheck directive at the beginning
		const newContent = `${TS_NOCHECK_DIRECTIVE}\n${content}`
		fs.writeFileSync(filePath, newContent, 'utf8')

		console.log(`Added @ts-nocheck to: ${filePath}`)
		return true
	} catch (error) {
		console.error(`Error processing ${filePath}:`, error.message)
		return false
	}
}

// Main function
function main() {
	// Orval injects generated files as arguments
	const files = process.argv.slice(2)

	console.log(`Script called with ${files.length} argument(s):`, files)

	let modifiedCount = 0

	// If no files provided, process the known Orval output files
	if (files.length === 0) {
		console.log('No files provided, processing known Orval output files...')
		const knownFiles = [
			'src/lib/api/index.ts',
			'src/lib/api/index.msw.ts',
			'src/lib/api/index.faker.ts',
		]

		for (const filePath of knownFiles) {
			if (fs.existsSync(filePath) && fs.statSync(filePath).isFile()) {
				if (addTsNocheckToFile(filePath)) {
					modifiedCount++
				}
			}
		}
	} else {
		// Process files provided as arguments
		for (const filePath of files) {
			if (fs.existsSync(filePath) && fs.statSync(filePath).isFile()) {
				if (addTsNocheckToFile(filePath)) {
					modifiedCount++
				}
			}
		}
	}

	if (modifiedCount > 0) {
		console.log(`✅ Added @ts-nocheck to ${modifiedCount} file(s)`)
	} else {
		console.log('ℹ️  No files needed @ts-nocheck directive')
	}
}

main()
